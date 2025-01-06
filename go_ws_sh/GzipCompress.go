package go_ws_sh

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

func GzipCompress(data []byte) ([]byte, bool, error) {
	var buf bytes.Buffer

	// 创建一个gzip.Writer，用于压缩数据
	gzWriter := gzip.NewWriter(&buf)
	defer func() {
		err := gzWriter.Close()
		if err != nil {
			fmt.Println("Error closing gzip writer:", err)
			return
		}
	}()
	// 将数据写入gzip.Writer
	_, err :=
		io.Copy(gzWriter, bytes.NewBuffer(data))
		// if err != nil && err != io.EOF {
		// 	fmt.Println("Error compressing data:", err)
		// 	return nil, true, err
		// }

		// 关闭gzip.Writer，确保所有数据都被压缩并写入缓冲区
	gzWriter.Flush()
	return buf.Bytes(), false, err
}
