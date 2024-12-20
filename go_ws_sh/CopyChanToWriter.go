package go_ws_sh

import "io"

func CopyChanToWriter(stdin io.WriteCloser, in_queue chan []byte) {
	for {
		var data, ok = <-in_queue
		if !ok {
			stdin.Close()
			return
		}
		stdin.Write(data)
	}

}
