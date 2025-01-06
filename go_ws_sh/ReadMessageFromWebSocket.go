package go_ws_sh

import (
	"log"

	"google.golang.org/protobuf/proto"
)

type WebsocketConnectionReadable interface {
	ReadMessage() (messageType int, p []byte, err error)
}

func ReadMessageFromWebSocket(c WebsocketConnectionReadable) (messageType int, p []byte, err error) {
	messageType, compressedData, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return messageType, nil, err
	}

	decompressedData, shouldReturn, err := GzipDeCompress(compressedData)
	if err != nil {
		log.Println("decompress:", err)
		return messageType, nil, err
	}

	if shouldReturn {
		return messageType, nil, err
	}
	var wsmsg = Wsmsg{}
	err = proto.Unmarshal(decompressedData, &wsmsg)
	if err != nil {
		log.Println("decompress:", err)
		return messageType, nil, err
	}
	return int(wsmsg.Type), wsmsg.Data, nil
}
