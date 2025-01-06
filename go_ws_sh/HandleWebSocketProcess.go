package go_ws_sh

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hertz-contrib/websocket"
	"github.com/linkedin/goavro/v2"
	"github.com/runletapp/go-console"
)

func SendTextMessage(conn *websocket.Conn, typestring string, body string, binaryandtextchannel chan WebsocketMessage) error {

	var data TextMessage
	data.Type = typestring
	data.Body = body
	databuf, err := json.Marshal(EncodeTextMessageToStringArray(data))
	if err != nil {
		return err
	}

	binaryandtextchannel <- WebsocketMessage{
		Body: databuf,
		Type: websocket.TextMessage,
	}

	return nil
}

type WebsocketMessage struct {
	Body []byte
	Type int
}

func HandleWebSocketProcess(session Session, codec *goavro.Codec, conn *websocket.Conn) error {
	var err error
	defer conn.WriteMessage(websocket.CloseMessage, []byte{})
	var binaryandtextchannel = make(chan WebsocketMessage)
	defer close(binaryandtextchannel)
	defer conn.Close()

	go func() {
		SendMessageToWebSocketLoop(conn, binaryandtextchannel)
	}()

	defer func() {
		defer conn.WriteMessage(websocket.CloseMessage, []byte{})

	}()

	var cmd console.Console = nil

	var mt int
	var message []byte

	mt, message, err = ReadMessageFromWebSocket(conn)
	var ok bool
	errclose, ok := err.(*websocket.CloseError)

	if ok {

		log.Println("close:", errclose)
		if cmd != nil {

			cmd.Kill()

		}

		return err
	} else if err != nil {
		log.Println("read:", err)
		return err
	}
	if mt == websocket.TextMessage {

		var array []any

		err = json.Unmarshal(message, &array)
		if err != nil {
			log.Println("read:", err)

			return err
		}
		log.Println("websocket recv text : ", (array))
		var data MessageSize
		err = DecodeMessageSizeFromStringArray(array, &data)
		if err != nil {
			log.Println("read:", err)

			return err
		}

		if data.Type == "resize" {
			log.Println("resize:", data.Cols, data.Rows)

			cmd, err = console.New(int(data.Cols), int(data.Rows))
			if err != nil {
				log.Println("resize:", err)
				return err
			}

		} else {
			log.Printf("ignored unknown recv text:%v", data)
			return errors.New("unknown recv text,first message console size expected")
		}

	} else {

		return errors.New("unknown recv binary,first message console size expected")
	}
	var Clear = func() {

		conn.Close()
		if cmd != nil {
			cmd.Kill()
		}

	}
	defer Clear()

	var stdin = cmd
	var stdout = cmd

	if err := cmd.Start(append([]string{session.Cmd}, session.Args...)); err != nil {
		log.Println(err)
		err := SendTextMessage(conn, "rejected", err.Error(), binaryandtextchannel)
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
	err = SendTextMessage(conn, "resolved", x, binaryandtextchannel)
	if err != nil {
		return err
	}

	go func() {

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

		for {

			var data, err = ReadFixedSizeFromReader(stdout, 1024*1024)
			if data == nil || nil != err {
				if err != nil {
					log.Println("encode:", err)
					return
				}
			}
			log.Println("server stdout received body:", data)

			var message = BinaryMessage{
				Type: "stdout",
				Body: data,
			}

			encoded, err := EncodeStructAvroBinary(codec, &message)

			if err != nil {
				log.Println("encode:", err)
				return
			}
			binaryandtextchannel <- WebsocketMessage{
				Body: encoded,
				Type: websocket.BinaryMessage,
			}

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
			log.Println("read:", err)
			break
		}
		if mt == websocket.TextMessage {

			var array []any

			err = json.Unmarshal(message, &array)
			if err != nil {
				log.Println("read:", err)

				return err
			}

			log.Println("websocket recv text : ", (array))
			var data MessageSize
			err = DecodeMessageSizeFromStringArray(array, &data)
			if err != nil {
				log.Println("read:", err)

				return err
			}

			if data.Type == "resize" {
				log.Println("resize:", data.Cols, data.Rows)
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
					log.Println("server stdin received body:", body)

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
