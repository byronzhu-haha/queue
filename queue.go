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
	initCap int
	realCap int
	head    int
	tail    int
	isFull  bool
	items   []interface{}
}

// NewQueue 如果cap <= 0, 则是容量无限大的队列
func NewQueue(cap int) Queue {
	var sli []interface{}
	if cap <= 0 {
		cap = -1
	} else {
		sli = make([]interface{}, cap)
	}
	return &queue{
		initCap: cap,
		realCap: cap,
		items:   sli,
	}
}

func (q *queue) isLimitSizeQueue() bool {
	return q.initCap == -1
}

func (q *queue) Cap() int {
	return q.realCap
}

func (q *queue) Len() int {
	if q.IsEmpty() {
		return 0
	}
	if q.IsFull() {
		return q.realCap
	}
	if q.tail > q.head {
		return q.tail - q.head
	}
	return q.realCap - q.head + q.tail
}

func (q *queue) IsEmpty() bool {
	return q.head == q.tail && !q.isFull
}

func (q *queue) IsFull() bool {
	return q.isFull
}

func (q *queue) Push(val interface{}) error {
	if q.isLimitSizeQueue() {
		return q.pushSlow(val)
	}
	if q.IsFull() {
		return ErrIsFull
	}
	q.items[q.tail] = val
	q.tail++
	if q.tail == q.realCap {
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
	q.realCap = cap(q.items)
	return nil
}

func (q *queue) Pull() (interface{}, error) {
	if q.IsEmpty() {
		return nil, ErrIsEmpty
	}
	if q.isLimitSizeQueue() {
		return q.pullSlow()
	}
	val := q.items[q.head]
	q.items[q.head] = nil
	q.head++
	if q.head == q.realCap {
		q.head = 0
	}
	q.isFull = false
	return val, nil
}

func (q *queue) pullSlow() (interface{}, error) {
	val := q.items[q.head]
	q.items = q.items[1:]
	q.tail--
	q.realCap = cap(q.items)
	return val, nil
}
