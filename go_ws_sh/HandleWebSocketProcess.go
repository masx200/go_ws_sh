package go_ws_sh

import (
	// "context"
	"encoding/json"
	"io"
	// "fmt"
	// "io"
	"log"
	"os/exec"

	"github.com/hertz-contrib/websocket"
	"github.com/linkedin/goavro/v2"
)

// SendTextMessage 通过WebSocket连接发送文本消息
// conn: WebSocket连接
// typestring: 消息类型
// body: 消息体
// mu: 互斥锁，用于同步写操作
// 返回错误，如果发送消息失败
func SendTextMessage(conn *websocket.Conn, typestring string, body string /*  mu *sync.Mutex */, binaryandtextchannel chan WebsocketMessage) error {

	var data TextMessage
	data.Type = typestring
	data.Body = body
	databuf, err := json.Marshal(EncodeTextMessageToStringArray(data))
	if err != nil {
		return err
	}

	// go func() {
	/* 这里不能开协程会乱序不可以 */
	binaryandtextchannel <- WebsocketMessage{
		Body: databuf,
		Type: websocket.TextMessage,
	}
	// }()
	//加一把锁在writemessage时使用,不能并发写入
	// mu.Lock()
	// defer mu.Unlock()
	// err = conn.WriteMessage(websocket.TextMessage, databuf)
	// if err != nil {
	// 	return fmt.Errorf("failed to send message: %w", err)
	// }

	return nil
}

type WebsocketMessage struct {
	Body []byte
	Type int
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
	var binaryandtextchannel = make(chan WebsocketMessage)
	defer close(binaryandtextchannel)
	defer conn.Close()
	// var in_queue = make(chan []byte)
	// var err_queue = make(chan []byte)
	// var out_queue = make(chan []byte)
	// defer close(out_queue)
	// defer close(err_queue)
	// defer close(in_queue)
	go func() {
		var err error
		for {
			//var encoded,ok <-  binaryandtextchannel
			encoded, ok := <-binaryandtextchannel
			if ok {
				// mubinary.Lock()
				// defer mubinary.Unlock()
				err = conn.WriteMessage(encoded.Type, encoded.Body)
				if err != nil {
					log.Println("write:", err)
					return
				}
			} else {
				break
			}
		}
	}()
	//加一把锁在writemessage时使用
	// var mutext sync.Mutex
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

	cmd := exec.Command(session.Cmd)

	var Clear = func() {
		//recover panic
		// defer func() {
		// 	if r := recover(); r != nil {
		// 		log.Printf("Recovered from panic: %v", r)
		// 	}
		// }()
		// close(out_queue)
		// close(err_queue)
		// close(in_queue)
		conn.Close()

		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}
	cmd.Args = session.Args

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Println(err)

		err := SendTextMessage(conn, "rejected", err.Error() /* &mutext */, binaryandtextchannel)
		if err != nil {
			return err
		}

		Clear()
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		err := SendTextMessage(conn, "rejected", err.Error() /*  &mutext */, binaryandtextchannel)
		if err != nil {
			return err
		}
		Clear()
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
		err := SendTextMessage(conn, "rejected", err.Error() /*  &mutext */, binaryandtextchannel)
		if err != nil {
			return err
		}
		Clear()
		return err
	}

	if err := cmd.Start(); err != nil {
		log.Println(err)
		err := SendTextMessage(conn, "rejected", err.Error() /* &mutext */, binaryandtextchannel)
		if err != nil {
			return err
		}
		Clear()
		return err
	}
	defer cmd.Process.Kill()
	x := "process " + session.Cmd + " started success"
	log.Println("resolved:" + x)
	err = SendTextMessage(conn, "resolved", x /* &mutext */, binaryandtextchannel)
	if err != nil {
		return err
	}
	// stdin.Write([]byte("ping qq.com" + "\n"))
	// go func() {
	// 	CopyReaderToChan(out_queue, stdout)
	// }()
	// go func() {
	// 	CopyReaderToChan(err_queue, stderr)
	// }()
	// go func() {
	// 	CopyChanToWriter(stdin, in_queue)

	// }()
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

			var data, err = ReadFixedSizeFromReader(stdout, 1024*1024)
			if data == nil || nil != err {
				if err != nil {
					log.Println("encode:", err)
					return
				}
			}
			// log.Printf("stdout recv Binary length: %v", len(data))
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
				return
			}
			binaryandtextchannel <- WebsocketMessage{
				Body: encoded,
				Type: websocket.BinaryMessage,
			}
			//encoded
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
			var data, err = ReadFixedSizeFromReader(stderr, 1024*1024)
			if data == nil || nil != err {
				if err != nil {
					log.Println("encode:", err)
					return
				}
			}
			// log.Printf("stderr recv Binary length: %v", len(data))
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
			binaryandtextchannel <- WebsocketMessage{
				Body: encoded,
				Type: websocket.BinaryMessage,
			} //encoded
			// mubinary.Lock()
			// defer mubinary.Unlock()
			// err = conn.WriteMessage(websocket.BinaryMessage, encoded)

			// if err != nil {
			// 	log.Println("write:", err)

			// }
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
			// log.Printf("websocket recv text length: %v", len(message))
			log.Printf("ignored recv text: %s", message)
		} else {
			// log.Printf("websocket recv Binary length: %v", len(message))

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
					// log.Println("server stdin received:", len(message))
					var body = md.Body
					// log.Println("body:", body)
					// go func() {
					//in_queue <- body
					stdin.Write(body)
					// }()
				} else {
					log.Println("ignored unknown type:", md.Type)
				}
				// }
			}
		}

	}
	return nil
}

// StreamReaderToChannel 将 io.Reader 中的数据流式地复制到一个字节切片通道中。
// 该函数会持续从 reader 读取数据，并将读取到的数据块发送到指定的通道 ch 中。
// 如果读取过程中发生错误，或者 reader 到达 EOF，函数将关闭通道 ch 并返回错误。
//
// 参数:
//   - reader: 数据源，实现了 io.Reader 接口。
//   - ch: 用于接收数据块的通道。
//
// 返回值:
//   - error: 如果读取过程中发生错误，返回该错误；否则返回 nil。
func ReadFixedSizeFromReader(stdin io.Reader, size int) ([]byte, error) {

	data := make([]byte, size)
	n, err := stdin.Read(data)
	if err != nil {
		// close(in_queue)
		return nil, err
	}
	// in_queue <- data[0:n]
	return data[0:n], nil

}
