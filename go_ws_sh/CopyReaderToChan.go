package go_ws_sh

import "io"

func CopyReaderToChan(in_queue chan []byte, stdin io.Reader) {

	for {
		data := make([]byte, 1024*1024)
		n, err := stdin.Read(data)
		if err != nil {
			close(in_queue)
			return
		}
		in_queue <- data[0:n]
	}
}
