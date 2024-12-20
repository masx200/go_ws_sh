package go_ws_sh

import (
	"fmt"
	// "github.com/cloudwego/hertz/pkg/app/client/retry"
)

type TextMessage struct {
	Type string `json:"type"`
	Body string `json:"body"`
}

// EncodeTextMessageToStringArray 将 TextMessage 结构体编码为 []string 数组
func EncodeTextMessageToStringArray(msg TextMessage) []string {
	return []string{msg.Type, msg.Body}
}

// DecodeTextMessageFromStringArray 将 []string 数组解码为 TextMessage 结构体
func DecodeTextMessageFromStringArray(strArray []string, result *TextMessage) error {
	if len(strArray) != 2 {
		return fmt.Errorf("invalid string array length, expected 2")
	}

	result.Type = strArray[0]

	result.Body = strArray[1]

	return nil
}
