package go_ws_sh

import (
	"log"

	goavro "github.com/linkedin/goavro/v2"
	"github.com/mitchellh/mapstructure"
)

type BinaryMessage struct {
	Type string `mapstructure:"type"`
	Body []byte `mapstructure:"body"`
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
func DecodeStructAvroBinary(codec *goavro.Codec, message []byte, result *any) error {
	decoded, undecoded, err := codec.NativeFromBinary(message)
	if len(undecoded) > 0 {
		log.Println("undecoded:", undecoded)
	}
	if err != nil {
		log.Println("decode:", err)
		return err
	}

	input := decoded
	err = mapstructure.Decode(input, &result)
	return err
}

// EncodeStructAvroBinary 将任意结构体编码为Avro二进制格式。
// 该函数首先将输入的结构体转换为map[any]interface{}类型，
// 然后使用提供的Avro编解码器（codec）将其编码为Avro二进制格式。
//
// 参数:
//   - codec: Avro编解码器，用于执行二进制编码。
//   - message: 指向要编码的结构体的指针。
//
// 返回值:
//   - []byte: 编码后的Avro二进制数据。
//   - error: 如果编码过程中发生错误，返回该错误。
func EncodeStructAvroBinary(codec *goavro.Codec, message *any) ([]byte, error) {
	var m map[any]interface{}
	err := mapstructure.Decode(message, &m)
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
