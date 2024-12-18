package go_ws_sh

import (
	"io"
	"slices"
	"sync"

	"github.com/gammazero/deque"
)

func init() {
	var _ io.Closer = NewQueue()
	var _ io.Reader = NewQueue()
	var _ io.Writer = NewQueue()
}

type Queue struct {
	data   *deque.Deque[[]byte]
	closed bool
	mu     *sync.Mutex
	cond   *sync.Cond
}

// Close implements io.Closer.
func (q *Queue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.closed = true
	q.cond.Broadcast()
	q.data = &deque.Deque[[]byte]{}
	return nil
}
func (q *Queue) Empty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.data.Len() == 0 && q.closed
}
func NewQueue() *Queue {
	var mu sync.Mutex
	x := sync.NewCond(&mu)
	return &Queue{data: &deque.Deque[[]byte]{},
		closed: false, cond: x, mu: &mu,
	}
}

func (q *Queue) Enqueue(value []byte) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.data.PushBack(value)
	q.cond.Signal()
}

func (q *Queue) Dequeue() []byte {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return nil
	}
	for q.data.Len() == 0 && !q.closed {
		q.cond.Wait() // Wait until there is data or the queue is closed
	}
	if q.data.Len() > 0 {
		value := q.data.Front()
		q.data.Remove(0)
		return value
	}
	return nil
}
func (q *Queue) EnqueueFront(value []byte) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.data.PushFront(value)
	q.cond.Signal()
}

// Read 从队列中读取数据到提供的字节切片p。
// 它作为队列的一个消费者，尝试从队列中取出数据。
// 如果队列已关闭，函数将返回0和EOF，表示没有更多数据可读。
// 如果队列中有数据，它将被取出并尽可能多地复制到p中。
// 如果p的大小小于队列中的下一条数据，剩余的数据将被重新入队到队列前端。
// 这个函数主要用于处理队列中的数据，将其消费或部分消费后重新入队。
// Read 从队列中读取数据到提供的字节切片 p。
// 参数:
//
//	p: 目标字节切片，用于存储从队列中读取的数据。
//
// 返回值:
//
//	n: 实际读取的字节数。
//	err: 错误信息，如果队列已关闭或没有数据可读，则返回 io.EOF。
func (q *Queue) Read(p []byte) (n int, err error) {
	if q.closed {
		return 0, io.EOF
	}
	value := q.Dequeue()
	if value == nil {
		return 0, io.EOF
	}
	if len(p) < len(value) {
		q.EnqueueFront(value[len(p):])
	}
	n = copy(p, value)

	return n, nil
}

func (q *Queue) Write(p []byte) (n int, err error) {
	q.Enqueue(slices.Clone(p))
	return len(p), nil
}
