package go_ws_sh

import (
	"io"
	"log"
)

func CopyChanToWriter(stdin io.WriteCloser, in_queue chan []byte) {
	for {
		var data, ok = <-in_queue
		if !ok {
			stdin.Close()
			return
		}
		x, y := stdin.Write(data)
		if y != nil {
			log.Println("CopyChanToWriter stdin recv Binary error: ", y)
			return
		}
		log.Println("CopyChanToWriter stdin recv Binary length: ", x)
	}

}
