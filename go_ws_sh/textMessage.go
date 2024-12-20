package go_ws_sh

import "fmt"

type TextMessage struct {
	Type string `json:"type"`
	Body string `json:"body"`
}

// EncodeTextMessageToStringArray 将 TextMessage 结构体编码为 []string 数组
func EncodeTextMessageToStringArray(msg *TextMessage) []string {
	return []string{msg.Type, msg.Body}
}

// DecodeTextMessageFromStringArray 将 []string 数组解码为 TextMessage 结构体
func DecodeTextMessageFromStringArray(strArray []string) (*TextMessage, error) {
	if len(strArray) != 2 {
		return &TextMessage{}, fmt.Errorf("invalid string array length, expected 2")
	}
	return &TextMessage{

		Type: strArray[0],

		Body: strArray[1],
	}, nil
}
