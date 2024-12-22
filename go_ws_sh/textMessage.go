package go_ws_sh

import (
	"fmt"
	// "github.com/cloudwego/hertz/pkg/app/client/retry"
)

type TextMessage struct {
	Type string
	Body string
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

type MessageSize struct {
	Type string
	Cols int
	Rows int
}

// EncodeTextMessageToStringArray 将 TextMessage 结构体编码为 []string 数组
func EncodeMessageSizeToStringArray(msg MessageSize) []any {
	return []any{msg.Type, msg.Cols, msg.Rows}
}

// DecodeTextMessageFromStringArray 将 []string 数组解码为 TextMessage 结构体
func DecodeMessageSizeFromStringArray(strArray []any, result *MessageSize) error {
	if len(strArray) != 3 {
		return fmt.Errorf("invalid string array length, expected 2")
	}
	var ok bool
	result.Type, ok = strArray[0].(string)

	if !ok {
		return fmt.Errorf("input is not a map[string]interface{} of MessageSize")
	}
	result.Cols, ok = strArray[1].(int)
	if !ok {
		return fmt.Errorf("input is not a map[string]interface{} of MessageSize")
	}
	result.Rows, ok = strArray[2].(int)
	if !ok {
		return fmt.Errorf("input is not a map[string]interface{} of MessageSize")
	}
	return nil
}
