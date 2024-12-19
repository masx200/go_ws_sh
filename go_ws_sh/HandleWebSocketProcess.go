package go_ws_sh

import (
	// "context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"

	"github.com/hertz-contrib/websocket"
	"github.com/linkedin/goavro/v2"
)

// SendTextMessage 通过WebSocket连接发送文本消息
// conn: WebSocket连接
// typestring: 消息类型
// body: 消息体
// mu: 互斥锁，用于同步写操作
// 返回错误，如果发送消息失败
func SendTextMessage(conn *websocket.Conn, typestring string, body string, mu *sync.Mutex) error {

	var data TextMessage
	data.Type = typestring
	data.Body = body
	databuf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	//加一把锁在writemessage时使用,不能并发写入
	mu.Lock()
	defer mu.Unlock()
	err = conn.WriteMessage(websocket.TextMessage, databuf)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

// HandleWebSocketProcess 处理WebSocket连接的整个生命周期。
// 该函数负责与客户端建立WebSocket连接，执行命令，并通过WebSocket发送和接收数据。
// 参数:
//
//	session: 包含要执行的命令和参数的会话信息。
//	codec: 用于编解码Avro消息的编解码器。
//	conn: 与客户端的WebSocket连接。
//
// 返回值:
//
//	如果执行过程中发生错误，则返回该错误。
func HandleWebSocketProcess(session Session, codec *goavro.Codec, conn *websocket.Conn) error {
	var binarychannel = make(chan []byte)
	defer close(binarychannel)
	//加一把锁在writemessage时使用
	var mutext sync.Mutex
	// defer mutext.Unlock()
	// var mubinary sync.Mutex
	defer func() {
		defer conn.WriteMessage(websocket.CloseMessage, []byte{})
		// if err := defer conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
		// 	log.Println(err)
		// }
	}()
	// var w, cancel = context.WithCancel(context.Background())
	// defer cancel()
	defer conn.Close()
	var in_queue = NewBlockingChannelDeque()
	var err_queue = NewBlockingChannelDeque()
	var out_queue = NewBlockingChannelDeque()
	defer out_queue.Close()
	defer err_queue.Close()
	defer in_queue.Close()
	cmd := exec.Command(session.Cmd)

	var Clear = func() {
		//recover panic
		// defer func() {
		// 	if r := recover(); r != nil {
		// 		log.Printf("Recovered from panic: %v", r)
		// 	}
		// }()
		out_queue.Close()
		err_queue.Close()
		in_queue.Close()
		conn.Close()

		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}
	cmd.Args = session.Args

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Println(err)

		err := SendTextMessage(conn, "rejected", err.Error(), &mutext)
		if err != nil {
			return err
		}

		Clear()
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		err := SendTextMessage(conn, "rejected", err.Error(), &mutext)
		if err != nil {
			return err
		}
		Clear()
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
		err := SendTextMessage(conn, "rejected", err.Error(), &mutext)
		if err != nil {
			return err
		}
		Clear()
		return err
	}

	if err := cmd.Start(); err != nil {
		log.Println(err)
		err := SendTextMessage(conn, "rejected", err.Error(), &mutext)
		if err != nil {
			return err
		}
		Clear()
		return err
	}
	defer cmd.Process.Kill()
	x := "process " + session.Cmd + " started success"
	log.Println(x)
	err = SendTextMessage(conn, "resolved", x, &mutext)
	if err != nil {
		return err
	}

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
		log.Println("process " + session.Cmd + " exit success")

		defer conn.WriteMessage(websocket.CloseMessage, []byte{})
		// cancel()
		// if err := defer conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
		// 	log.Println(err)
		// }

		// conn.WriteControl()
		Clear()
		conn.Close()

	}()
	go func() {

		for {

			data := out_queue.Dequeue()
			if data == nil {
				break
			}
			log.Printf("stdout recv Binary length: %v", len(data))
			var message = BinaryMessage{
				Type: "stdout",
				Body: data,
			}

			encoded, err := EncodeStructAvroBinary(codec, &message)
			// encoded, err := codec.BinaryFromNative(nil, map[string]interface{}{
			// 	"type": "stdout",
			// 	"body": data,
			// })
			if err != nil {
				log.Println("encode:", err)
				continue
			}
			binarychannel <- encoded
			// mubinary.Lock()
			// defer mubinary.Unlock()
			// err = conn.WriteMessage(websocket.BinaryMessage, encoded)

			// if err != nil {
			// 	log.Println("write:", err)

			// }
		}
	}()
	go func() {

		for {
			data := err_queue.Dequeue()
			if data == nil {
				break
			}
			log.Printf("stderr recv Binary length: %v", len(data))
			var message = BinaryMessage{
				Type: "stderr",
				Body: data,
			}

			encoded, err := EncodeStructAvroBinary(codec, &message)
			// encoded, err := codec.BinaryFromNative(nil, map[string]interface{}{
			// 	"type": "stderr",
			// 	"body": data,
			// })
			if err != nil {
				log.Println("encode:", err)
				continue
			}
			binarychannel <- encoded
			// mubinary.Lock()
			// defer mubinary.Unlock()
			// err = conn.WriteMessage(websocket.BinaryMessage, encoded)

			// if err != nil {
			// 	log.Println("write:", err)

			// }
		}
	}()
	go func() {

		for {
			//var encoded,ok <-  binarychannel
			encoded, ok := <-binarychannel
			if ok {
				// mubinary.Lock()
				// defer mubinary.Unlock()
				err = conn.WriteMessage(websocket.BinaryMessage, encoded)
				if err != nil {
					log.Println("write:", err)
				}
			} else {
				break
			}
		}
	}()
	for {

		// select {
		// case <-w.Done():
		// 	log.Println("exit done conn close")
		// 	return nil
		// default:

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
			log.Printf("websocket recv text length: %v", len(message))
			log.Printf("ignored recv text: %s", message)
		} else {
			log.Printf("websocket recv Binary length: %v", len(message))

			var result BinaryMessage
			undecoded, err := DecodeStructAvroBinary(codec, message, &result)
			if len(undecoded) > 0 {

				log.Println("undecoded:", undecoded)

			}
			if err != nil {
				log.Println("decode:", err)

			} else {
				// log.Printf("recv binary: %s", decoded)
				var md = result
				if md.Type == "stdin" {
					var body = md.Body
					in_queue.Enqueue(body)
				} else {
					log.Println("ignored unknown type:", md.Type)
				}
				// }
			}
		}

	}
	return nil
}
