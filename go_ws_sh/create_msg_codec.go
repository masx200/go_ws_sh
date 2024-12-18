package go_ws_sh

import goavro "github.com/linkedin/goavro/v2"

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
