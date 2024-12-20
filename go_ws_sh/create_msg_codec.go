package go_ws_sh

import (
	"fmt"
	"log"

	goavro "github.com/linkedin/goavro/v2"
	// "github.com/mitchellh/mapstructure"
)

type BinaryMessage struct {
	Type string
	Body []byte
}

// create_msg_codec 创建一个用于编码和解码消息的Apache Avro编解码器。
// 该函数定义了一个消息的Avro模式，该模式包括一个字符串类型的"type"字段和一个字节类型的"body"字段。
// 返回值是一个*goavro.Codec类型的编解码器实例，以及一个错误值，如果创建编解码器时出现错误，则该错误值会被设置。
func create_msg_codec() (*goavro.Codec, error) {
	const schemaJSON = `
{
    "type": "record",
    "name": "message",
    "fields": [
        {
            "name": "type",
            "type": "string"
        },
        {
            "name": "body",
            "type": "bytes"
        }
    ]
}
`
	codec, err := goavro.NewCodec(schemaJSON)
	return codec, err
}

// DecodeStructAvroBinary 解析Avro二进制消息到指定的结构体。
// 参数:
//
//	codec: Avro编解码器，用于解析二进制消息。
//	message: 待解析的二进制消息。
//	result: 解析后的数据将被存储的结构体指针。
//
// 返回值:
//
//	如果解析过程中发生错误，则返回错误。
func DecodeStructAvroBinary(codec *goavro.Codec, message []byte, result *BinaryMessage) ([]byte, error) {
	decoded, undecoded, err := codec.NativeFromBinary(message)
	if len(undecoded) > 0 {
		log.Println("undecoded:", undecoded)
	}
	if err != nil {
		log.Println("decode:", err)
		return undecoded, err
	}

	input := decoded
	err = DecodeBinaryMessage(input, result)
	return undecoded, err
}

func DecodeBinaryMessage(input interface{}, result *BinaryMessage) error {
	if result == nil {
		return fmt.Errorf("result is nil")
	}

	// 断言 input 为 map[string]interface{}
	inputMap, ok := input.(map[string]interface{})
	if !ok {
		return fmt.Errorf("input is not a map[string]interface{}")
	}

	// 获取 "type" 字段并赋值给 result.Type
	typeValue, ok := inputMap["type"].(string)
	if !ok {
		return fmt.Errorf("type field is not a string")
	}
	result.Type = typeValue

	// 获取 "body" 字段并赋值给 result.Body
	bodyValue, ok := inputMap["body"].([]byte)
	if !ok {
		return fmt.Errorf("body field is not a []byte")
	}
	result.Body = bodyValue

	return nil
}

// EncodeStructAvroBinary 将任意结构体编码为Avro二进制格式。
// 该函数首先将输入的结构体转换为map[string]interface{}类型，
// 然后使用提供的Avro编解码器（codec）将其编码为Avro二进制格式。
//
// 参数:
//   - codec: Avro编解码器，用于执行二进制编码。
//   - message: 指向要编码的结构体的指针。
//
// 返回值:
//   - []byte: 编码后的Avro二进制数据。
//   - error: 如果编码过程中发生错误，返回该错误。
func EncodeStructAvroBinary(codec *goavro.Codec, message *BinaryMessage) ([]byte, error) {
	var m map[string]interface{} = make(map[string]interface{})
	err := EncodeBinaryMessage(message, m)
	if err != nil {
		log.Println("decode:", err)
		return nil, err
	}
	encoded, err := codec.BinaryFromNative(nil, m)
	if err != nil {
		log.Println("encode:", err)
		return nil, err
	}
	return encoded, nil

}

func EncodeBinaryMessage(message *BinaryMessage, m map[string]interface{}) error {
	if message == nil {
		return fmt.Errorf("message is nil")
	}
	if nil == m {

		return fmt.Errorf("map is nil")
	}
	// *m = map[string]interface{}{
	// 	"type": message.Type,
	// 	"body": message.Body,
	// }
	m["type"] = message.Type
	m["body"] = message.Body
	return nil
}
