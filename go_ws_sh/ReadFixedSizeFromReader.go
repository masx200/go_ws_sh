package go_ws_sh

import "io"

// StreamReaderToChannel 将 io.Reader 中的数据流式地复制到一个字节切片通道中。
// 该函数会持续从 reader 读取数据，并将读取到的数据块发送到指定的通道 ch 中。
// 如果读取过程中发生错误，或者 reader 到达 EOF，函数将关闭通道 ch 并返回错误。
//
// 参数:
//   - reader: 数据源，实现了 io.Reader 接口。
//   - ch: 用于接收数据块的通道。
//
// 返回值:
//   - error: 如果读取过程中发生错误，返回该错误；否则返回 nil。
func ReadFixedSizeFromReader(stdin io.Reader, size int) ([]byte, error) {

	data := make([]byte, size)
	n, err := stdin.Read(data)
	if err != nil {
		// close(in_queue)
		return nil, err
	}
	// in_queue <- data[0:n]
	return data[0:n], nil

}
