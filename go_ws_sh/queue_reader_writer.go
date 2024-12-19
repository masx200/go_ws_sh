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

// BlockingChannelDeque 是一个阻塞通道双端队列。
// 它使用双端队列作为底层数据结构，并提供了线程安全的操作。
// 该结构用于在生产者和消费者之间高效地传递数据，同时保持数据的顺序。
type BlockingChannelDeque struct {
	// data 是一个双端队列，存储着字节切片。它是BlockingChannelDeque的核心组件，
	// 负责实际数据的存储和操作。
	data *deque.Deque[[]byte]

	// closed 表示BlockingChannelDeque是否已关闭。当closed为true时，表示不能再向
	// BlockingChannelDeque中添加数据，但仍然可以移除已有的数据直到队列为空。
	closed bool

	// mu 是一个互斥锁，用于保护BlockingChannelDeque的成员变量，确保在同一时刻
	// 只有一个线程可以修改BlockingChannelDeque的状态。
	mu *sync.Mutex

	// cond 是一个条件变量，与互斥锁mu一起使用。当BlockingChannelDeque为空或满时，
	// 线程可以等待cond以获取信号继续操作，从而实现阻塞和唤醒机制。
	cond *sync.Cond
}

// Done implements BlockingChannel.
// Done is a method of BlockingChannelDeque used to wait for the deque to close.
// This method is primarily used to block the caller until the BlockingChannelDeque is closed.
// It should be noted that this method will not return immediately if the deque has not been closed,
// but will wait until the deque is closed.
func (q *BlockingChannelDeque) Done() {
	// Acquire the lock to ensure thread safety during the operation.
	q.mu.Lock()
	defer q.mu.Unlock()

	// If the deque has already been closed, return the method directly.
	if q.closed {
		return
	}

	// Wait in a loop until the deque is closed.
	// The use of conditional variables here allows the current goroutine to wait until notified,
	// thus avoiding unnecessary busy waiting and improving program efficiency.
	for !q.closed {
		q.cond.Wait()
	}

	//return
}

// Closed implements BlockingDeque.
// Closed 判断通道是否已关闭。
// 该方法通过检查BlockingChannelDeque实例的closed属性来确定通道的关闭状态。
// 使用互斥锁确保线程安全，防止多个协程同时修改或读取closed状态。
func (q *BlockingChannelDeque) Closed() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.closed
}

// IsEmpty implements BlockingDeque.
// IsEmpty 检查阻塞通道双端队列是否为空。
// 该方法通过检查内部数据结构的长度来确定队列是否为空。
// 使用互斥锁确保在检查长度时队列不会被修改。
func (q *BlockingChannelDeque) IsEmpty() bool {
	// 加锁以确保线程安全
	q.mu.Lock()
	defer q.mu.Unlock()
	// 返回队列是否为空的布尔值
	return q.data.Len() == 0
}

// PushFront implements BlockingDeque.

// PushBack implements BlockingDeque.
// PushBack 将一个字节切片添加到BlockingChannelDeque的末尾。
// 该方法在内部使用锁来确保线程安全，并在队列关闭时防止添加新元素。
// 当队列已满或需要通知等待的协程时，它还会触发条件变量。
//
// 参数:
//
//	item []byte - 要添加到队列末尾的字节切片。
//
// 返回值:
//
//	error - 如果队列已关闭，则返回io.EOF；否则返回nil。
func (q *BlockingChannelDeque) PushBack(item []byte) error {
	// 上锁以确保线程安全。
	q.mu.Lock()
	defer q.mu.Unlock()

	// 检查队列是否已关闭，如果是，则返回EOF错误。
	if q.closed {
		return io.EOF
	}

	// 将项添加到队列末尾。
	q.data.PushBack(item)

	// 触发条件变量以通知可能在等待的协程。
	q.cond.Signal()

	// 成功添加项后返回nil。
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
