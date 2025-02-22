package go_ws_sh

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hertz-contrib/websocket"
	"github.com/linkedin/goavro/v2"
)

func SendTextMessage(conn *websocket.Conn, typestring string, body string, binaryandtextchannel *SafeChannel[WebsocketMessage]) error {

	var data TextMessage
	data.Type = typestring
	data.Body = body
	databuf, err := json.Marshal(EncodeTextMessageToStringArray(data))
	if err != nil {
		return err
	}

	binaryandtextchannel.Send(WebsocketMessage{
		Body: databuf,
		Type: websocket.TextMessage,
	},
	)
	return nil
}

type WebsocketMessage struct {
	Body []byte
	Type int
}

func HandleWebSocketProcess(session Session, codec *goavro.Codec, conn *websocket.Conn) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in panic", r)
		}
	}()
	var err2 error
	defer conn.WriteMessage(websocket.CloseMessage, []byte{})
	var binaryandtextchannel = NewSafeChannel[WebsocketMessage]()
	defer (binaryandtextchannel).Close()
	defer conn.Close()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in panic", r)
			}
		}()
		SendMessageToWebSocketLoop(conn, binaryandtextchannel)
	}()

	defer func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in panic", r)
			}
		}()
		defer conn.WriteMessage(websocket.CloseMessage, []byte{})

	}()

	cmd, err2 := handleWebSocketConnection(conn)
	if err2 != nil {
		return sendErrorMessageToWebSocket(conn, err2)
	}
	var Clear = func() {

		conn.Close()
		if cmd != nil {
			cmd.Kill()
		}

	}

	if cmd == nil {
		return errors.New("cmd is nil")
	}
	defer Clear()
	var stdin = cmd
	var stdout = cmd
	cmd.SetCWD(session.Dir)
	if err := cmd.Start(append([]string{session.Cmd}, session.Args...)); err != nil {
		log.Println(err)
		err := sendErrorMessageToWebSocket(conn, err)
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
	err2 = SendTextMessage(conn, "resolved", x, binaryandtextchannel)
	if err2 != nil {
		return err2
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in panic", r)
			}
		}()
		state, err := cmd.Wait()
		if err != nil {
			log.Println(err)
			Clear()
			return
		}
		log.Println("process " + session.Cmd + " exit success" + " code:" + fmt.Sprintf("%d", state.ExitCode()))

		defer conn.WriteMessage(websocket.CloseMessage, []byte{})

		Clear()
		conn.Close()

	}()
	go func() {
		//panic: send on closed channel
		//recover
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in panic", r)
			}
		}()
		for {

			var data, err = ReadFixedSizeFromReader(stdout, 1024*1024)
			if data == nil || nil != err {
				if err != nil {
					log.Println("encode:", err)
					return
				}
			}
			// log.Println("server stdout received body:", data)

			var message = BinaryMessage{
				Type: "stdout",
				Body: data,
			}

			encoded, err := EncodeStructAvroBinary(codec, &message)

			if err != nil {
				log.Println("encode:", err)
				return
			}
			binaryandtextchannel.Send(WebsocketMessage{
				Body: encoded,
				Type: websocket.BinaryMessage,
			})

		}
	}()

	for {

		mt, message, err := ReadMessageFromWebSocket(conn)
		if errclose, ok := err.(*websocket.CloseError); ok {

			log.Println("close:", errclose)
			if cmd != nil {

				cmd.Kill()

			}
			break
		}
		if err != nil {
			log.Println("read4:", err)
			break
		}
		if mt == websocket.TextMessage {

			var array []any

			err = json.Unmarshal(message, &array)
			if err != nil {
				log.Println("read5:", err)

				return err
			}

			// log.Println("websocket recv text : ", (array))
			var data MessageSize
			err = DecodeMessageSizeFromStringArray(array, &data)
			if err != nil {
				log.Println("read6:", err)

				return err
			}

			if data.Type == "resize" {
				// log.Println("resize:", data.Cols, data.Rows)
				if cmd != nil {
					cmd.SetSize(int(data.Cols), int(data.Rows))
				}

			} else {
				log.Printf("ignored unknown recv text:%v", data)
				return errors.New("ignored unknown recv text console message size expected")
			}

		} else {

			var result BinaryMessage
			undecoded, err := DecodeStructAvroBinary(codec, message, &result)
			if len(undecoded) > 0 {

				log.Println("undecoded:", undecoded)

			}
			if err != nil {
				log.Println("decode:", err)

			} else {

				var md = result
				if md.Type == "stdin" {

					var body = md.Body
					// log.Println("server stdin received body:", body)

					stdin.Write(body)

				} else {
					log.Println("ignored unknown type:", md.Type)
					return errors.New("ignored unknown type  stdin expected ")

				}

			}
		}

	}
	return nil
}
