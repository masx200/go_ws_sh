package go_ws_sh

import (
	"io"
	"log"
	"os/exec"

	"github.com/hertz-contrib/websocket"
	"github.com/linkedin/goavro/v2"
)

func HandleWebSocketProcess(session Session, codec *goavro.Codec, conn *websocket.Conn) {

	defer conn.Close()
	var in_queue = NewQueue()
	var err_queue = NewQueue()
	var out_queue = NewQueue()
	defer out_queue.Close()
	defer err_queue.Close()
	defer in_queue.Close()
	cmd := exec.Command(session.Cmd)

	var Clear = func() {
		//recover panic
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}

		conn.Close()
		out_queue.Close()
		err_queue.Close()
		in_queue.Close()
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}
	cmd.Args = session.Args

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Println(err)
		Clear()
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		Clear()
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
		Clear()
		return
	}

	if err := cmd.Start(); err != nil {
		log.Println(err)
		Clear()
		return
	}
	defer cmd.Process.Kill()

	go func() {
		io.Copy(out_queue, stdout)
	}()
	go func() {
		io.Copy(err_queue, stderr)
	}()
	go func() {
		io.Copy(stdin, in_queue)

	}()
	go func() {
		if err := cmd.Wait(); err != nil {
			log.Println(err)
			Clear()
			return
		}

		conn.Close()
	}()
	go func() {

		for {
			data := out_queue.Dequeue()
			if data == nil {
				break
			}

			encoded, err := codec.BinaryFromNative(nil, map[string]interface{}{
				"type": "stdout",
				"body": data,
			})
			if err != nil {
				log.Println("encode:", err)
				continue
			}

			err = conn.WriteMessage(websocket.BinaryMessage, encoded)

			if err != nil {
				log.Println("write:", err)

			}
		}
	}()
	go func() {

		for {
			data := err_queue.Dequeue()
			if data == nil {
				break
			}

			encoded, err := codec.BinaryFromNative(nil, map[string]interface{}{
				"type": "stderr",
				"body": data,
			})
			if err != nil {
				log.Println("encode:", err)
				continue
			}

			err = conn.WriteMessage(websocket.BinaryMessage, encoded)

			if err != nil {
				log.Println("write:", err)

			}
		}
	}()

	for {
		mt, message, err := conn.ReadMessage()
		if err, ok := err.(*websocket.CloseError); ok {

			log.Println("close:", err)

			cmd.Process.Kill()
			break
		}
		if err != nil {
			log.Println("read:", err)
			break
		}
		if mt == websocket.TextMessage {
			log.Printf("ignored recv text: %s", message)
		} else {
			decoded, undecoded, err := codec.NativeFromBinary(message)
			if len(undecoded) > 0 {

				log.Println("undecoded:", undecoded)

			}
			if err != nil {
				log.Println("decode:", err)

			} else {
				// log.Printf("recv binary: %s", decoded)
				var md = decoded.(map[string]interface{})
				if md["type"] == "stdin" {
					var body = md["body"].([]byte)
					in_queue.Enqueue(body)
				} else {
					log.Println("ignored unknown type:", md["type"])
				}
			}

		}

	}
}
