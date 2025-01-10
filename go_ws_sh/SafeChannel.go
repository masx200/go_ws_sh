package go_ws_sh

import (
	"fmt"
	"sync"
)

// SafeChannel 是一个包装了 channel 的结构体，提供了安全的发送、接收和关闭操作。
type SafeChannel[T any] struct {
	ch     chan T
	closed bool
	mu     sync.Mutex
	// wg     sync.WaitGroup
}

// NewSafeChannel 创建一个新的 SafeChannel 实例。
func NewSafeChannel[T any](buffer ...int) *SafeChannel[T] {
	if len(buffer) == 0 {
		buffer = []int{0}
	}
	sc := &SafeChannel[T]{
		ch: make(chan T, buffer[0]),
	}
	// sc.wg.Add(1)
	// go sc.monitor()
	return sc
}
func (sc *SafeChannel[T]) IsClosed() bool {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	return sc.closed
}

// Send 尝试向 channel 发送数据，如果 channel 已经关闭则返回 false。
func (sc *SafeChannel[T]) Send(v T) bool {
	//panic: send on closed channel
	//recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	sc.mu.Lock()

	if sc.closed {
		sc.mu.Unlock()
		return false
	}
	sc.mu.Unlock()
	// select {
	//case
	sc.ch <- v
	return true
	// return true
	// default:
	// 如果 channel 已满，可以选择等待或立即返回。
	// 这里选择立即返回 false。
	// return false
	//}
}

// Receive 尝试从 channel 接收数据，如果 channel 已关闭并且没有更多数据，则返回零值和 false。
func (sc *SafeChannel[T]) Receive() (T, bool) {
	sc.mu.Lock()

	if sc.closed && len(sc.ch) == 0 {
		var zero T // 零值
		sc.mu.Unlock()
		return zero, false
	}
	sc.mu.Unlock()
	v, ok := <-sc.ch
	return v, ok
}

// Close 安全地关闭 channel。
func (sc *SafeChannel[T]) Close() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if !sc.closed {
		close(sc.ch)
		sc.closed = true
	}
}
