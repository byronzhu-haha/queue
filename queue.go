package queue

import "errors"

type Queue interface {
	Cap() int
	Len() int
	IsEmpty() bool
	IsFull() bool
	Push(val interface{}) error
	Pull() (interface{}, error)
}

var (
	ErrIsFull  = errors.New("queue is full")
	ErrIsEmpty = errors.New("queue is empty")
)

type queue struct {
	cap        int
	head       int
	tail       int
	isAutoGrow bool
	isFull     bool
	items      []interface{}
}

// NewQueue 如果cap <= 0, 则是容量无限大的队列
func NewQueue(cap int) Queue {
	return newQueue(cap)
}

func newQueue(cap int) *queue {
	var (
		sli        []interface{}
		isAutoGrow bool
	)
	if cap <= 0 {
		cap = 0
		isAutoGrow = true
	} else {
		sli = make([]interface{}, cap)
	}
	return &queue{
		cap:        cap,
		isAutoGrow: isAutoGrow,
		items:      sli,
	}
}

func (q *queue) Cap() int {
	return q.cap
}

func (q *queue) Len() int {
	if q.IsEmpty() {
		return 0
	}
	if q.IsFull() {
		return q.cap
	}
	if q.tail > q.head {
		return q.tail - q.head
	}
	return q.cap - q.head + q.tail
}

func (q *queue) IsEmpty() bool {
	return q.head == q.tail && !q.isFull
}

func (q *queue) IsFull() bool {
	return q.isFull
}

func (q *queue) Push(val interface{}) error {
	if q.isAutoGrow {
		return q.pushSlow(val)
	}
	if q.IsFull() {
		return ErrIsFull
	}
	q.items[q.tail] = val
	q.tail++
	if q.tail == q.cap {
		q.tail = 0
	}
	if q.head == q.tail {
		q.isFull = true
	}
	return nil
}

func (q *queue) pushSlow(val interface{}) error {
	q.items = append(q.items, val)
	q.tail++
	q.cap = cap(q.items)
	return nil
}

func (q *queue) Pull() (interface{}, error) {
	if q.IsEmpty() {
		return nil, ErrIsEmpty
	}
	if q.isAutoGrow {
		return q.pullSlow()
	}
	val := q.items[q.head]
	q.items[q.head] = nil
	q.head++
	if q.head == q.cap {
		q.head = 0
	}
	q.isFull = false
	return val, nil
}

func (q *queue) pullSlow() (interface{}, error) {
	nc := q.reduceCap()
	val := q.items[q.head]
	q.items = q.items[1:]
	q.tail--
	q.cap = nc
	return val, nil
}

func (q *queue) reduceCap() (newCap int) {
	n, oc := q.Len(), q.Cap()
	if oc > 31 && n <= (oc>>2) {
		newCap = oc >> 1
		nq := make([]interface{}, n, newCap)
		copy(nq, q.items)
		q.items = nq
	} else {
		newCap = oc
	}
	return newCap
}

func (q *queue) swap(i, j int) {
	q.items[i], q.items[j] = q.items[j], q.items[i]
}

func (q *queue) pop() (interface{}, error) {
	if q.IsEmpty() {
		return nil, ErrIsEmpty
	}
	if q.isAutoGrow {
		return q.popSlow()
	}
	idx := q.tail-1
	val := q.items[idx]
	q.items[idx] = nil
	q.tail--
	if q.tail <= 0 {
		q.tail = 0
	}
	q.isFull = false
	return val, nil
}

func (q *queue) popSlow() (interface{}, error) {
	nc := q.reduceCap()
	idx := q.tail-1
	val := q.items[idx]
	q.items = q.items[:idx]
	q.tail--
	q.cap = nc
	return val, nil
}
