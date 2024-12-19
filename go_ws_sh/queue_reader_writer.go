package go_ws_sh

import (
	"io"
	"slices"
	"sync"

	"github.com/gammazero/deque"
)

// init函数用于初始化BlockingChannelDeque，确保其符合多个接口的要求。
// 这里的目的是验证BlockingChannelDeque实现了io.Closer、io.Reader、io.Writer、BlockingDeque和BlockingChannel等接口。
// 通过这种方式，可以在程序启动时对BlockingChannelDeque的接口实现进行检查，确保其行为符合预期。
func init() {
	// 将NewBlockingChannelDeque的返回值赋值给io.Closer接口变量，验证其是否实现了Close方法。
	var _ io.Closer = NewBlockingChannelDeque()
	// 将NewBlockingChannelDeque的返回值赋值给io.Reader接口变量，验证其是否实现了Read方法。
	var _ io.Reader = NewBlockingChannelDeque()
	// 将NewBlockingChannelDeque的返回值赋值给io.Writer接口变量，验证其是否实现了Write方法。
	var _ io.Writer = NewBlockingChannelDeque()
	// 将NewBlockingChannelDeque的返回值赋值给BlockingDeque接口变量，验证其是否实现了特定的BlockingDeque方法。
	var _ BlockingDeque = NewBlockingChannelDeque()
	// 将NewBlockingChannelDeque的返回值赋值给BlockingChannel接口变量，验证其是否实现了特定的BlockingChannel方法。
	var _ BlockingChannel = NewBlockingChannelDeque()
}

type BlockingChannelDeque struct {
	data   *deque.Deque[[]byte]
	closed bool
	mu     *sync.Mutex
	cond   *sync.Cond
}

// Done implements BlockingChannel.
func (q *BlockingChannelDeque) Done() {

	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return
	}
	for !q.closed {
		q.cond.Wait()
	}

}

// Closed implements BlockingDeque.
func (q *BlockingChannelDeque) Closed() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.closed
}

// IsEmpty implements BlockingDeque.
func (q *BlockingChannelDeque) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.data.Len() == 0
}

// PushFront implements BlockingDeque.

// PushBack implements BlockingDeque.
func (q *BlockingChannelDeque) PushBack(item []byte) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return io.EOF
	}
	q.data.PushBack(item)
	q.cond.Signal()
	return nil
}

// Size implements BlockingDeque.
func (q *BlockingChannelDeque) Size() int {

	return q.data.Len()
}

// TakeFirst implements BlockingDeque.
func (q *BlockingChannelDeque) TakeFirst() ([]byte, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	x := q.Dequeue()
	if x == nil {
		return nil, false
	}
	return x, true
}

// TakeLast implements BlockingDeque.
func (q *BlockingChannelDeque) TakeLast() ([]byte, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return nil, false
	}
	for q.data.Len() == 0 && !q.closed {
		q.cond.Wait() // Wait until there is data or the queue is closed
	}
	if q.data.Len() == 0 {
		return nil, false
	}
	x := q.data.Back()
	q.data.Remove(q.data.Len() - 1)
	if x == nil {
		return nil, false
	}
	return x, true
}

// Close implements io.Closer.
func (q *BlockingChannelDeque) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.closed = true
	q.cond.Broadcast()
	q.data.Clear()
	return nil
}
func (q *BlockingChannelDeque) Empty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.data.Len() == 0 && q.closed
}
func NewBlockingChannelDeque() *BlockingChannelDeque {
	var mu sync.Mutex
	x := sync.NewCond(&mu)
	return &BlockingChannelDeque{data: &deque.Deque[[]byte]{},
		closed: false, cond: x, mu: &mu,
	}
}

func (q *BlockingChannelDeque) Enqueue(value []byte) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return io.EOF
	}
	q.data.PushBack(value)
	q.cond.Signal()
	return nil
}

func (q *BlockingChannelDeque) Dequeue() []byte {
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
func (q *BlockingChannelDeque) PushFront(value []byte) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return io.EOF
	}
	q.data.PushFront(value)
	q.cond.Signal()
	return nil
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
func (q *BlockingChannelDeque) Read(p []byte) (n int, err error) {
	if q.closed {
		return 0, io.EOF
	}
	value := q.Dequeue()
	if value == nil {
		return 0, io.EOF
	}
	if len(p) < len(value) {
		q.PushFront(value[len(p):])
	}
	n = copy(p, value)

	return n, nil
}

func (q *BlockingChannelDeque) Write(p []byte) (n int, err error) {
	if q.closed {
		return 0, io.EOF
	}
	err = q.Enqueue(slices.Clone(p))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

type BlockingDeque interface {
	// 在队尾插入元素，如果队列已满，则阻塞等待
	PushBack(item []byte) error
	// 在队首插入元素，如果队列已满，则阻塞等待
	PushFront(item []byte) error
	// 从队尾移除元素，如果队列为空，则阻塞等待
	TakeLast() ([]byte, bool)
	// 从队首移除元素，如果队列为空，则阻塞等待
	TakeFirst() ([]byte, bool)
	// 检查队列是否为空
	IsEmpty() bool
	// 获取队列的大小
	Size() int
	Closed() bool
}
type BlockingChannel interface {
	// 在队尾插入元素，如果队列已满，则阻塞等待
	Enqueue(item []byte) error
	// 在队首插入元素，如果队列已满，则阻塞等待
	// PushFront(item []byte) error
	// 从队尾移除元素，如果队列为空，则阻塞等待
	// TakeLast() ([]byte, bool)
	// 从队首移除元素，如果队列为空，则阻塞等待
	Dequeue() []byte
	// 检查队列是否为空
	IsEmpty() bool
	// 获取队列的大小
	Size() int
	Closed() bool
	Done()
}
