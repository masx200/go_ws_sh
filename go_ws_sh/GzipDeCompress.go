package go_ws_sh

import (
	"bytes"
	"io"

	"github.com/klauspost/pgzip"
)

func GzipDeCompress(b []byte) ([]byte, error) {
	reader, err := pgzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	var buf *bytes.Buffer = &bytes.Buffer{}
	_, err = io.Copy(buf, reader)
	if err != nil {
		return nil, err
	}
	reader.Close()
	// 读取解压缩后的数据
	var decompressed = buf.Bytes()
	return decompressed, nil
}
