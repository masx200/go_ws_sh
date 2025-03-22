package go_ws_sh

import (
	"bytes"
	"log"

	"io"

	"github.com/klauspost/pgzip"
)

func GzipCompress(data []byte) ([]byte, error) {
	var buf *bytes.Buffer = &bytes.Buffer{}

	// 创建一个gzip.Writer，用于压缩数据
	gzWriter := pgzip.NewWriter(buf)
	defer func() {
		err := gzWriter.Close()
		if err != nil {
			log.Println("Error closing gzip writer:", err)
			return
		}
	}()
	// 将数据写入gzip.Writer
	_, err :=
		io.Copy(gzWriter, bytes.NewBuffer(data))
		// if err != nil && err != io.EOF {
		// 	log.Println("Error compressing data:", err)
		// 	return nil, true, err
		// }

		// 关闭gzip.Writer，确保所有数据都被压缩并写入缓冲区
	gzWriter.Flush()
	gzWriter.Close()
	return buf.Bytes(), err
}
