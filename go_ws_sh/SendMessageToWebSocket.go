package go_ws_sh

import (
	"log"

	"github.com/hertz-contrib/websocket"
	"google.golang.org/protobuf/proto"
)

type WebsocketConnectionWritableClosable interface {
	WriteMessage(messageType int, data []byte) error
	Close() error
}

func SendMessageToWebSocket(conn WebsocketConnectionWritableClosable, encoded WebsocketMessage) error {
	var err error

	var b []byte
	var wsmsg = Wsmsg{}
	wsmsg.Type = int32(encoded.Type)
	wsmsg.Data = encoded.Body
	b, err = proto.Marshal(&wsmsg)
	if err != nil {
		log.Println("write:", err)
		return err
	}
	bg, shouldReturn, err := GzipCompress(b)
	if err != nil {
		log.Println("write:", err)
		return err
	}
	if shouldReturn {
		return err
	}
	err = conn.WriteMessage(websocket.BinaryMessage, bg)
	return err
}
func SendMessageToWebSocketLoop(conn WebsocketConnectionWritableClosable, binaryandtextchannel chan WebsocketMessage) {
	defer conn.Close()

	for {
		var err error
		encoded, ok := <-binaryandtextchannel
		if ok {
			err = SendMessageToWebSocket(conn, encoded)
			if err != nil {
				log.Println("write:", err)
				return
			}
		} else {
			break
		}
	}

}
