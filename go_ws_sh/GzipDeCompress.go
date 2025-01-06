package go_ws_sh

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func GzipDeCompress(b []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// 读取解压缩后的数据
	decompressed, err := ioutil.ReadAll(reader)
	if err != nil {
		return decompressed, err
	}

	return decompressed, nil
}
