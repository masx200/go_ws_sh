package go_ws_sh

import (
	"fmt"
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
	Cols int64
	Rows int64
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

	cols, ok := strArray[1].(float64)
	if !ok {
		return fmt.Errorf("input is not a map[string]interface{} of MessageSize")
	}
	result.Cols = int64(cols)

	rows, ok := strArray[2].(float64)
	if !ok {
		return fmt.Errorf("input is not a map[string]interface{} of MessageSize")
	}
	result.Rows = int64(rows)
	return nil
}
