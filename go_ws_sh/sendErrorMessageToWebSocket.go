package go_ws_sh

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/hertz-contrib/websocket"
	"github.com/runletapp/go-console"
	"google.golang.org/protobuf/proto"
)

func sendErrorMessageToWebSocket(conn *websocket.Conn, err2 error) error {
	var err error
	var body = err2.Error()
	var typestring = "rejected"
	var data TextMessage
	data.Type = typestring
	data.Body = body
	databuf, err := json.Marshal(EncodeTextMessageToStringArray(data))
	if err != nil {
		return err
	}
	var encoded = WebsocketMessage{
		Body: databuf,
		Type: websocket.TextMessage,
	}
	var b []byte
	var wsmsg = Wsmsg{}
	wsmsg.Type = int32(encoded.Type)
	wsmsg.Data = encoded.Body
	b, err = proto.Marshal(&wsmsg)
	if err != nil {
		log.Println("write3:", err)
		return err
	}
	bg, err := GzipCompress(b)
	if err != nil {
		log.Println("write4:", err)
		return err
	}

	// log.Printf("SendMessageToWebSocket before encode %v %v\n", encoded.Type, encoded.Body)
	// log.Printf("SendMessageToWebSocket after encode %v %v\n", websocket.BinaryMessage, bg)
	err = conn.WriteMessage(websocket.BinaryMessage, bg)
	if err != nil {
		return err
	}
	return err2
}

func handleWebSocketConnection(conn *websocket.Conn) (console.Console, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in panic", r)
		}
	}()
	var cmd console.Console = nil
	var err error
	var mt int
	var message []byte

	mt, message, err = ReadMessageFromWebSocket(conn)
	// log.Printf("first message %v %v %v \n", mt, message, err)
	var ok bool
	errclose, ok := err.(*websocket.CloseError)

	if ok {

		log.Println("close:", errclose)
		// if cmd != nil {

		// 	cmd.Kill()

		// }

		return nil, err
	} else if err != nil {
		log.Println("read1:", err)
		err2 := sendErrorMessageToWebSocket(conn, errors.New("unknown recv message,first message console size expected"))
		if err2 != nil {
			return nil, err2
		}
		return nil, err
	}
	if mt == websocket.TextMessage {

		var array []any

		err = json.Unmarshal(message, &array)
		if err != nil {
			log.Println("read2:", err)
			err2 := sendErrorMessageToWebSocket(conn, errors.New("unknown recv text,first message console size expected"))
			if err2 != nil {
				return nil, err2
			}
			return nil, err
		}
		// log.Println("websocket recv text : ", (array))
		var data MessageSize
		err = DecodeMessageSizeFromStringArray(array, &data)
		if err != nil {
			log.Println("read3:", err)
			err2 := sendErrorMessageToWebSocket(conn, errors.New("unknown recv text,first message console size expected"))
			if err2 != nil {
				return nil, err2
			}
			return nil, err
		}

		if data.Type == "resize" {
			// log.Println("resize:", data.Cols, data.Rows)

			cmd, err = console.New(int(data.Cols), int(data.Rows))
			if err != nil {
				log.Println("resize:", err)
				err2 := sendErrorMessageToWebSocket(conn, errors.New("unknown recv text,first message console size expected"))
				if err2 != nil {
					return nil, err2
				}
				return nil, err
			}
			return cmd, nil
		} else {
			log.Printf("ignored unknown recv text:%v", data)
			err2 := sendErrorMessageToWebSocket(conn, errors.New("unknown recv text,first message console size expected"))
			if err2 != nil {
				return nil, err2
			}
			return nil, errors.New("unknown recv text,first message console size expected")
		}

	} else {
		log.Printf("ignored unknown recv binary:%v", message)
		err2 := sendErrorMessageToWebSocket(conn, errors.New("unknown recv binary,first message console size expected"))
		if err2 != nil {
			return nil, err2
		}
		return nil, errors.New("unknown recv binary,first message console size expected")
	}

}
