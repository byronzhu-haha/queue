package queue

import (
	"errors"
	"fmt"
)

type priorityQueue struct {
	queue    *queue
	lessFunc func(i, j int) bool
}

type Item struct {
	idx      int
	Priority int
	Value    interface{}
}

const defaultCap = 1024

var (
	ErrValueIsNotItemType = errors.New("value's type is not Item")
	ErrValueIsNil         = errors.New("value is nil")
	ErrQueueIsNil = errors.New("queue is nil")
	ErrQueueIsNotPQ = errors.New("queue is not priority queue")
	ErrQueueCantSetLessFunc = errors.New("queue can't set less func because it is not empty")
)

// NewPriorityQueueWithCap 创建指定容量的队列，该队列不会自动触发扩容操作
func NewPriorityQueueWithCap(cap int) Queue {
	if cap <= 0 {
		cap = defaultCap
	}
	return &priorityQueue{queue: newQueue(cap)}
}

// NewNewPriorityQueue 创建一个理论上容量无限大的队列，会自动触发扩容操作
func NewNewPriorityQueue() Queue {
	return &priorityQueue{queue: newQueue(0)}
}

// Ascend 根据优先级权重设置升序，必须在队列创建后插入元素前使用，即队列为空的情况下使用
func Ascend(queue Queue) error {
	pq, err := isValidQueue(queue)
	if err != nil {
		return err
	}
	if !pq.IsEmpty() {
		return ErrQueueCantSetLessFunc
	}
	pq.ascend()
	return nil
}

// Ascend 根据优先级权重设置降序，必须在队列创建后插入元素前使用，即队列为空的情况下使用
func Descend(queue Queue) error {
	pq, err := isValidQueue(queue)
	if err != nil {
		return err
	}
	if !pq.IsEmpty() {
		return ErrQueueCantSetLessFunc
	}
	pq.descend()
	return nil
}

func isValidQueue(queue Queue) (*priorityQueue, error) {
	if queue == nil {
		return nil, ErrQueueIsNil
	}
	pq, ok := queue.(*priorityQueue)
	if !ok {
		return nil, ErrQueueIsNotPQ
	}
	if pq == nil {
		return nil, ErrQueueIsNil
	}
	return pq, nil
}

func (pq *priorityQueue) Len() int {
	return pq.queue.Len()
}

func (pq *priorityQueue) Cap() int {
	return pq.queue.Cap()
}

func (pq *priorityQueue) IsEmpty() bool {
	return pq.queue.IsEmpty()
}

func (pq *priorityQueue) IsFull() bool {
	return pq.queue.IsFull()
}

func (pq *priorityQueue) Push(val interface{}) error {
	value, ok := val.(*Item)
	if !ok {
		return ErrValueIsNotItemType
	}
	if value == nil {
		return ErrValueIsNil
	}

	value.idx = pq.queue.Len() - 1
	err := pq.queue.Push(value)
	if err != nil {
		return err
	}
	pq.up(value.idx)

	return nil
}

func (pq *priorityQueue) Pull() (interface{}, error) {
	if pq.IsEmpty() {
		return nil, ErrIsEmpty
	}
	end := pq.Len()-1
	pq.swap(0, end)
	pq.down(0, end)
	val, err := pq.queue.pop()
	if err != nil {
		return nil, err
	}
	val.(*Item).idx = -1
	return val, nil
}

func (pq *priorityQueue) get(idx int) *Item {
	return pq.queue.items[idx].(*Item)
}

func (pq *priorityQueue) less(i, j int) bool {
	if pq.lessFunc != nil {
		return pq.lessFunc(i, j)
	}
	return pq.get(i).Priority > pq.get(j).Priority
}

func (pq *priorityQueue) ascend() {
	pq.lessFunc = func(i, j int) bool {
		return pq.get(i).Priority < pq.get(j).Priority
	}
}

func (pq *priorityQueue) descend() {
	pq.lessFunc = func(i, j int) bool {
		return pq.get(i).Priority > pq.get(j).Priority
	}
}

func (pq *priorityQueue) swap(i, j int) {
	pq.queue.swap(i, j)
	pq.get(i).idx = j
	pq.get(j).idx = i
}

// up 上浮
func (pq *priorityQueue) up(idx int) {
	for {
		parent := (idx - 1) / 2
		if parent == idx || !pq.less(idx, parent) {
			break
		}
		pq.swap(parent, idx)
		idx = parent
	}
}

// down 下沉
func (pq *priorityQueue) down(start, end int) bool {
	i := start
	for {
		child := 2*i + 1
		if child >= end || child < 0 {
			break
		}

		lChild := child
		if rChild := child + 1; rChild < end && pq.less(rChild, child) {
			lChild = rChild
		}
		if !pq.less(lChild, i) {
			break
		}
		pq.swap(i, lChild)
		i = lChild
	}
	return i > start
}

func (i *Item) String() string {
	return fmt.Sprintf("{p: %d, v: %+v, idx: %d}", i.Priority, i.Value, i.idx)
}
