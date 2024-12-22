package go_ws_sh

import (
	// "context"
	"encoding/json"
	"errors"
	"fmt"
	// "io"
	// "fmt"
	// "io"
	"log"
	// "os/exec"

	"github.com/hertz-contrib/websocket"
	"github.com/linkedin/goavro/v2"
	"github.com/runletapp/go-console"
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
	var err error
	defer conn.WriteMessage(websocket.CloseMessage, []byte{})
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

	var cmd console.Console = nil

	/* 读取第一条消息,
	获取终端大小,
	否则断开连接 */
	//exec.Command(session.Cmd, session.Args...)
	var mt int
	var message []byte

	mt, message, err = conn.ReadMessage()
	var ok bool
	err, ok = err.(*websocket.CloseError)

	if ok {

		log.Println("close:", err)
		if cmd != nil {
			// if cmd.Process != nil {
			cmd.Kill()
			// }
		}
		// break
		return err
	} else /* if err != nil
	 */
	//  {
	// log.Println("read:", err)
	// return err
	// }
	if mt == websocket.TextMessage {
		// log.Printf("websocket recv text length: %v", len(message))
		// log.Printf("ignored recv text: %s", message)
		var array []any
		//parse json data

		err = json.Unmarshal(message, &array)
		if err != nil {
			log.Println("read:", err)
			//return
			// log.Printf("ignored recv text: %s", message)
			return err
		}
		var data MessageSize
		err = DecodeMessageSizeFromStringArray(array, &data)
		if err != nil {
			log.Println("read:", err)
			//return
			// log.Printf("ignored recv text: %s", message)
			return err
		}
		// log.Println("websocket recv text length: ", len(message))
		if data.Type == "resize" {
			log.Println("resize:", data.Cols, data.Rows)
			if cmd != nil {
				cmd.SetSize(data.Cols, data.Rows)
			} else {
				cmd, err = console.New(data.Cols, data.Rows)
				if err != nil {
					log.Println("resize:", err)
					return err
				}
			}
			// defer os.Exit(0)
			// return
			//break
		} else {
			log.Printf("ignored unknown recv text:%v", data)
			return errors.New("unknown recv text,first message console size expected")
		}
		/* else if data.Type == "resolved" {
			log.Println("resolved:", data.Body)
		} */

	} else {

		return errors.New("unknown recv binary,first message console size expected")
	}
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
		if cmd != nil {
			cmd.Kill()
		}

	}
	defer Clear()
	// cmd.Args = session.Args

	// stdin, err := cmd.StdinPipe()
	// if err != nil {
	// 	log.Println(err)

	// 	err := SendTextMessage(conn, "rejected", err.Error() /* &mutext */, binaryandtextchannel)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	Clear()
	// 	return err
	// }
	var stdin = cmd
	var stdout = cmd
	// var stderr = cmd
	// stdout, err := cmd.StdoutPipe()
	// if err != nil {
	// 	log.Println(err)
	// 	err := SendTextMessage(conn, "rejected", err.Error() /*  &mutext */, binaryandtextchannel)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	Clear()
	// 	return err
	// }
	// stderr, err := cmd.StderrPipe()
	// if err != nil {
	// 	log.Println(err)
	// 	err := SendTextMessage(conn, "rejected", err.Error() /*  &mutext */, binaryandtextchannel)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	Clear()
	// 	return err
	// }

	if err := cmd.Start(append([]string{session.Cmd}, session.Args...)); err != nil {
		log.Println(err)
		err := SendTextMessage(conn, "rejected", err.Error() /* &mutext */, binaryandtextchannel)
		if err != nil {
			return err
		}
		Clear()
		return err
	}
	defer func() {
		if cmd != nil {
			cmd.Kill()
		}
	}()
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

		state, err := cmd.Wait()
		if err != nil {
			log.Println(err)
			Clear()
			return
		}
		log.Println("process " + session.Cmd + " exit success" + " code:" + fmt.Sprintf("%d", state.ExitCode()))

		defer conn.WriteMessage(websocket.CloseMessage, []byte{})
		// cancel()
		// if err := defer conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
		// 	log.Println(err)
		// }

		// conn.WriteControl()
		Clear()
		conn.Close()
		// return
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
			log.Println("server stdout received body:", data)
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
	// go func() {

	// 	for {
	// 		var data, err = ReadFixedSizeFromReader(stderr, 1024*1024)
	// 		if data == nil || nil != err {
	// 			if err != nil {
	// 				log.Println("encode:", err)
	// 				return
	// 			}
	// 		}
	// 		log.Println("server stderr received body:", data)
	// 		// log.Printf("stderr recv Binary length: %v", len(data))
	// 		var message = BinaryMessage{
	// 			Type: "stderr",
	// 			Body: data,
	// 		}

	// 		encoded, err := EncodeStructAvroBinary(codec, &message)
	// 		// encoded, err := codec.BinaryFromNative(nil, map[string]interface{}{
	// 		// 	"type": "stderr",
	// 		// 	"body": data,
	// 		// })
	// 		if err != nil {
	// 			log.Println("encode:", err)
	// 			continue
	// 		}
	// 		binaryandtextchannel <- WebsocketMessage{
	// 			Body: encoded,
	// 			Type: websocket.BinaryMessage,
	// 		} //encoded
	// 		// mubinary.Lock()
	// 		// defer mubinary.Unlock()
	// 		// err = conn.WriteMessage(websocket.BinaryMessage, encoded)

	// 		// if err != nil {
	// 		// 	log.Println("write:", err)

	// 		// }
	// 	}
	// }()

	for {

		// select {
		// case <-w.Done():
		// 	log.Println("exit done conn close")
		// 	return nil
		// default:

		mt, message, err := conn.ReadMessage()
		if err, ok := err.(*websocket.CloseError); ok {

			log.Println("close:", err)
			if cmd != nil {
				// if cmd.Process != nil {
				cmd.Kill()
				// }
			}
			break
		}
		if err != nil {
			log.Println("read:", err)
			break
		}
		if mt == websocket.TextMessage {
			// log.Printf("websocket recv text length: %v", len(message))
			// log.Printf("ignored recv text: %s", message)
			var array []any
			//parse json data

			err = json.Unmarshal(message, &array)
			if err != nil {
				log.Println("read:", err)
				//return
				// log.Printf("ignored recv text: %s", message)
				return err
			}
			var data MessageSize
			err = DecodeMessageSizeFromStringArray(array, &data)
			if err != nil {
				log.Println("read:", err)
				//return
				// log.Printf("ignored recv text: %s", message)
				return err
			}
			// log.Println("websocket recv text length: ", len(message))
			if data.Type == "resize" {
				log.Println("resize:", data.Cols, data.Rows)
				if cmd != nil {
					cmd.SetSize(data.Cols, data.Rows)
				}
				// defer os.Exit(0)
				// return
				//break
			} else {
				log.Printf("ignored unknown recv text:%v", data)
			}
			/* else if data.Type == "resolved" {
				log.Println("resolved:", data.Body)
			} */

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
					log.Println("server stdin received body:", body)
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
