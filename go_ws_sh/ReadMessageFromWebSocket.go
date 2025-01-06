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
		log.Println("read10:", err)
		return messageType, nil, err
	}
	// log.Printf("ReadMessageFromWebSocket before decode %v %v \n", messageType, compressedData)
	decompressedData, err := GzipDeCompress(compressedData)
	if err != nil {
		log.Println("decompress1:", err)
		return messageType, decompressedData, err
	}

	var wsmsg = Wsmsg{}
	err = proto.Unmarshal(decompressedData, &wsmsg)
	if err != nil {
		log.Println("decompress2:", err)
		return messageType, decompressedData, err
	}
	// log.Printf("ReadMessageFromWebSocket after decode %v %v \n", wsmsg.Type, wsmsg.Data)
	return int(wsmsg.Type), wsmsg.Data, nil
}
