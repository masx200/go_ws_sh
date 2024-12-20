package go_ws_sh

import (
	"io"
	"math"
	// "slices"
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

// type IterableQueue interface {
// 	Iterator() IteratorQueue
// }
// type IteratorQueue interface {
// 	Next() ([]byte, bool)
// }
// type IteratorQueueImplementation struct {
// 	q *BlockingChannelDeque
// 	p int
// }

// BlockingChannelDeque 是一个阻塞通道双端队列。
// 它使用双端队列作为底层数据结构，并提供了线程安全的操作。
// 该结构用于在生产者和消费者之间高效地传递数据，同时保持数据的顺序。
type BlockingChannelDeque struct {
	// data 是一个双端队列，存储着字节切片。它是BlockingChannelDeque的核心组件，
	// 负责实际数据的存储和操作。
	data *deque.Deque[byte]

	// closed 表示BlockingChannelDeque是否已关闭。当closed为true时，表示不能再向
	// BlockingChannelDeque中添加数据，但仍然可以移除已有的数据直到队列为空。
	closed bool

	// mu 是一个互斥锁，用于保护BlockingChannelDeque的成员变量，确保在同一时刻
	// 只有一个线程可以修改BlockingChannelDeque的状态。
	mu   *sync.Mutex
	read *sync.Mutex
	// cond 是一个条件变量，与互斥锁mu一起使用。当BlockingChannelDeque为空或满时，
	// 线程可以等待cond以获取信号继续操作，从而实现阻塞和唤醒机制。
	cond *sync.Cond
}

// PeekFirst implements BlockingDeque.
func (q *BlockingChannelDeque) PeekFirst() (byte, bool) {
	q.read.Lock()
	defer q.read.Unlock()
	if q.closed {
		return 0, false
	}
	for (q.data.Len() == 0) && !q.closed {
		q.cond.Wait()
	}
	if q.closed {
		return 0, false
	}
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.data.Front(), true
	// q.mu.Lock()
	// defer q.mu.Unlock()

	// if q.closed {
	// 	return nil, false
	// }
	// for q.data.Len() == 0 && !q.closed {
	// 	q.cond.Wait() // Wait until there is data or the queue is closed
	// }
	// if q.data.Len() > 0 {
	// 	value := q.data.Front()
	// 	// q.data.Remove(0)
	// 	return value, true
	// }
	// return nil, false
}

// PeekLast 从 BlockingChannelDeque 中获取最后一个元素而不移除它。
//
// 参数：
//   - q: 指向 BlockingChannelDeque 实例的指针。
//
// 返回值：
//   - byte: 队列中的最后一个字节。如果队列已关闭或为空，则返回 0。
//   - bool: 表示操作是否成功。如果队列已关闭或为空，则返回 false。
//
// 该函数通过条件变量等待，直到队列中有元素或队列被关闭，
// 并使用互斥锁确保线程安全。
// PeekLast implements BlockingDeque.
func (q *BlockingChannelDeque) PeekLast() (byte, bool) {
	q.read.Lock()
	defer q.read.Unlock()
	if q.closed {
		return 0, false
	}
	for (q.data.Len() == 0) && !q.closed {
		q.cond.Wait()
	}
	if q.closed {
		return 0, false
	}
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.data.Back(), true

	// q.mu.Lock()
	// defer q.mu.Unlock()
	// if q.closed {
	// 	return nil, false
	// }
	// for q.data.Len() == 0 && !q.closed {
	// 	q.cond.Wait() // Wait until there is data or the queue is closed
	// }
	// if q.data.Len() == 0 {
	// 	return nil, false
	// }
	// x := q.data.Back()
	// // q.data.Remove(q.data.Len() - 1)
	// if x == nil {
	// 	return nil, false
	// }
	// return x, true
}

// Wait implements BlockingChannel.
// Wait is a method of BlockingChannelDeque used to wait for the deque to close.
// This method is primarily used to block the caller until the BlockingChannelDeque is closed.
// It should be noted that this method will not return immediately if the deque has not been closed,
// but will wait until the deque is closed.
func (q *BlockingChannelDeque) Wait() {
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
func (q *BlockingChannelDeque) PushBack(item byte) error {
	// 上锁以确保线程安全。
	q.mu.Lock()
	defer q.mu.Unlock()

	// 检查队列是否已关闭，如果是，则返回EOF错误。
	if q.closed {
		return io.EOF
	}

	// 将项添加到队列末尾。
	q.data.PushBack((item))

	// 触发条件变量以通知可能在等待的协程。
	q.cond.Signal()

	// 成功添加项后返回nil。
	return nil
}

// Size implements BlockingDeque.
// Size 返回BlockingChannelDeque中的元素数量。
// 该方法通过锁定互斥锁来确保线程安全，在解锁后返回元素数量。
// 注意：使用互斥锁是因为BlockingChannelDeque可能在多个线程间共享，
// 锁定可以防止在计算大小时对deque进行修改，从而确保数据一致性。
func (q *BlockingChannelDeque) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.data.Len()
}

// TakeFirst implements BlockingDeque.
// TakeFirst 从BlockingChannelDeque中取出并返回第一个元素。
// 该方法首先锁定互斥锁以确保线程安全，然后尝试从队列尾部取出元素。
// 如果队列为空，即没有元素可取，则返回nil和false。
// 如果成功取出元素，则返回该元素和true。
// 注意：该方法假设调用者已经处理了潜在的nil指针问题。
func (q *BlockingChannelDeque) TakeFirst() (byte, bool) {
	q.read.Lock()
	defer q.read.Unlock()
	if q.closed {
		return 0, false
	}
	for (q.data.Len() == 0) && !q.closed {
		q.cond.Wait()
	}
	if q.closed {
		return 0, false
	}
	q.mu.Lock()
	defer q.mu.Unlock()

	x := q.data.Front()
	q.data.Remove(0)
	return x, true
	// q.mu.Lock()
	// defer q.mu.Unlock()

	// if q.closed {
	// 	return nil, false
	// }
	// for q.data.Len() == 0 && !q.closed {
	// 	q.cond.Wait() // Wait until there is data or the queue is closed
	// }
	// if q.data.Len() > 0 {
	// 	value := q.data.Front()
	// 	q.data.Remove(0)
	// 	return value, true
	// }
	// return nil, false
}

// TakeLast implements BlockingDeque.
// TakeLast removes and returns the last element from the queue.
// If the queue is closed or empty, it returns nil and false.
// This method blocks until the queue has elements or is closed.
// It is thread-safe.
func (q *BlockingChannelDeque) TakeLast() (byte, bool) {
	q.read.Lock()
	defer q.read.Unlock()
	if q.closed {
		return 0, false
	}
	for (q.data.Len() == 0) && !q.closed {
		q.cond.Wait()
	}
	if q.closed {
		return 0, false
	}
	q.mu.Lock()
	defer q.mu.Unlock()

	x := q.data.Back()
	q.data.Remove(q.data.Len() - 1)
	return x, true

	// q.mu.Lock()
	// defer q.mu.Unlock()
	// if q.closed {
	// 	return nil, false
	// }
	// for q.data.Len() == 0 && !q.closed {
	// 	q.cond.Wait() // Wait until there is data or the queue is closed
	// }
	// if q.data.Len() == 0 {
	// 	return nil, false
	// }
	// x := q.data.Back()
	// q.data.Remove(q.data.Len() - 1)
	// if x == nil {
	// 	return nil, false
	// }
	// return x, true
}

// Close implements io.Closer.
// Close 关闭阻塞通道队列，并清理所有资源。
//
// 该方法首先获取互斥锁以确保线程安全，然后标记队列为已关闭状态，并唤醒所有等待的协程。
// 最后，清除队列中的所有数据。此方法不接受任何参数，也不返回任何错误。
func (q *BlockingChannelDeque) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.closed = true
	q.cond.Broadcast()
	q.data.Clear()
	return nil
}

// Empty 判断阻塞通道队列是否为空。
// 该方法通过检查内部数据结构的长度来确定队列是否为空。
// 使用互斥锁确保线程安全，防止多个goroutine同时修改队列状态。
func (q *BlockingChannelDeque) Empty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.data.Len() == 0
}

// NewBlockingChannelDeque 创建并返回一个新的BlockingChannelDeque实例。
// 该函数使用了互斥锁（mu）和条件变量（cond）来同步对双端队列（deque）的访问，
// 以实现线程安全。通过初始化一个关闭标志（closed）来标记队列是否已关闭，
// 确保在多线程环境下对队列的操作是安全和有序的。
func NewBlockingChannelDeque() *BlockingChannelDeque {
	// 初始化互斥锁，用于保护对队列的访问。
	var mu sync.Mutex
	// 创建一个新的条件变量，用于线程间的通信，确保对队列的操作是线程安全的。
	x := sync.NewCond(&mu)
	// 返回一个新的BlockingChannelDeque实例，其中包含一个空的双端队列、一个未关闭的状态、
	// 一个条件变量和一个互斥锁的引用。
	return &BlockingChannelDeque{
		data:   &deque.Deque[byte]{},
		closed: false,
		cond:   x,
		mu:     &mu,
		read:   &sync.Mutex{},
	}
}

// Enqueue 将一个字节切片值添加到BlockingChannelDeque的队列中。
// 该方法在内部使用互斥锁来确保线程安全，并在队列关闭时阻止进一步的入队操作。
// 当入队成功时，它会通知可能在等待队列非空的消费者。
// 参数:
//
//	value - 要入队的字节切片值。
//
// 返回值:
//
//	如果队列已经关闭，则返回io.EOF错误，否则返回nil。
func (q *BlockingChannelDeque) Enqueue(value []byte) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return io.EOF
	}
	if len(value) == 0 {
		return nil
	}
	q.data.Grow(len(value))
	for _, b := range value {
		q.data.PushBack((b))
	}
	// q.data.PushBack((value))
	q.cond.Signal()
	return nil
}

// Dequeue dequeues a message from the BlockingChannelDeque.
// This function is blocking; if there are no messages in the queue, it will wait until a message is available or the queue is closed.
// Returns nil if the queue is closed and there are no messages available.
func (q *BlockingChannelDeque) Dequeue() []byte {
	q.read.Lock()
	defer q.read.Unlock()
	if q.closed {
		return nil
	}
	for (q.data.Len() == 0) && !q.closed {
		q.cond.Wait()
	}
	if q.closed {
		return nil
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	var p = make([]byte, q.data.Len())
	var minsize = q.data.Len()

	for i := 0; i < int(minsize); i++ {
		p[i] = q.data.At(i)

	}
	for i := 0; i < int(minsize); i++ {
		q.data.Remove(0)
	}
	return p
	// q.mu.Lock()
	// defer q.mu.Unlock()

	// if q.closed {
	// 	return nil
	// }
	// for q.data.Len() == 0 && !q.closed {
	// 	q.cond.Wait() // Wait until there is data or the queue is closed
	// }
	// if q.data.Len() > 0 {
	// 	value := q.data.Front()
	// 	q.data.Remove(0)
	// 	return value
	// }
	// return nil
}

// PushFront 将一个字节切片作为消息添加到队列的前端。
// 此方法在添加消息前会检查队列是否已关闭，如果关闭则返回错误。
// 参数:
//
//	value - 要添加到队列的消息，类型为 []byte。
//
// 返回值:
//
//	如果队列已关闭，返回 io.EOF 错误；否则返回 nil。
func (q *BlockingChannelDeque) PushFront(value byte) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return io.EOF
	}
	q.data.PushFront((value))
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
	q.read.Lock()
	defer q.read.Unlock()
	if q.closed {
		return 0, io.EOF
	}
	for (q.data.Len() == 0) && !q.closed {
		q.cond.Wait()
	}
	if q.closed {
		return 0, io.EOF
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	var minsize = int(math.Min(float64(len(p)), float64(q.data.Len())))

	for i := 0; i < int(minsize); i++ {
		p[i] = q.data.At(i)

	}
	for i := 0; i < int(minsize); i++ {
		q.data.Remove(0)
	}
	return minsize, nil
	// //先尝试从队列中取出数据
	// //PeekFirst()
	// // first, ok := q.PeekFirst()

	// // if !ok {
	// // 	return 0, io.EOF
	// // }
	// // if len(first) < len(p) {

	// // }
	// value := q.Dequeue()
	// if value == nil {
	// 	return 0, io.EOF
	// }
	// if len(p) < len(value) {
	// 	err = q.PushFront((value[len(p):]))
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// }
	// n = copy(p, value)

	// return n, nil
}

// Write 向BlockingChannelDeque写入数据。
// 该方法在BlockingChannelDeque关闭时返回EOF错误。
// 参数:
//
//	p []byte: 要写入的数据。
//
// 返回值:
//
//	n int: 成功写入的字节数。
//	err error: 如果写入失败，返回错误。
func (q *BlockingChannelDeque) Write(p []byte) (n int, err error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return n, io.EOF
	}
	if len(p) == 0 {
		return n, nil
	}
	q.data.Grow(len(p))
	for _, b := range p {
		q.data.PushBack((b))
	}
	// q.data.PushBack((value))
	n = len(p)
	q.cond.Signal()
	return n, nil
}

// BlockingDeque 是一个阻塞双端队列的接口，提供了在队列两端进行插入和移除操作的能力
type BlockingDeque interface {
	PeekLast() (byte, bool)
	// TakeFirst 从队首移除元素，如果队列为空，则阻塞等待直到队列中有元素可用
	PeekFirst() (byte, bool)
	// PushBack 在队尾插入元素，如果队列已满，则阻塞等待直到队列有空余位置
	PushBack(item byte) error
	// PushFront 在队首插入元素，如果队列已满，则阻塞等待直到队列有空余位置
	PushFront(item byte) error
	// TakeLast 从队尾移除元素，如果队列为空，则阻塞等待直到队列中有元素可用
	TakeLast() (byte, bool)
	// TakeFirst 从队首移除元素，如果队列为空，则阻塞等待直到队列中有元素可用
	TakeFirst() (byte, bool)
	// IsEmpty 检查队列是否为空
	IsEmpty() bool
	// Size 获取队列的大小
	Size() int
	// Closed 检查队列是否已经关闭
	Closed() bool
}

// BlockingChannel 定义了一个阻塞通道接口，用于在队列满或空时阻塞操作
type BlockingChannel interface {
	PushBack(item byte) error
	// 在队尾插入元素，如果队列已满，则阻塞等待
	Enqueue(item []byte) error
	// 在队首插入元素，如果队列已满，则阻塞等待
	PushFront(item byte) error
	// 从队尾移除元素，如果队列为空，则阻塞等待
	TakeLast() (byte, bool)
	TakeFirst() (byte, bool)
	// 从队首移除元素，如果队列为空，则阻塞等待
	Dequeue() []byte
	// 检查队列是否为空
	IsEmpty() bool
	// 获取队列的大小
	Size() int
	// 检查队列是否已关闭
	Closed() bool
	// 完成队列的操作，通常用于通知其他协程可以停止等待
	Wait()
}
