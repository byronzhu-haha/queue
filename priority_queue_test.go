package queue

import (
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	pq := NewNewPriorityQueue()
	err := pq.Push("abc")
	equalVal(t, ErrValueIsNotItemType, err, "push error")

	var it *Item
	err = pq.Push(it)
	equalVal(t, ErrValueIsNil, err, "push error")

	it = &Item{
		idx:      0,
		Priority: 7,
		Value:    "666",
	}
	err = pq.Push(it)
	equalVal(t, nil, err, "push error")
	equalVal(t, 1, pq.Len(), "len error")
	equalVal(t, pq.(*priorityQueue).queue.cap, pq.Cap(), "cap error")
	equalVal(t, false, pq.IsEmpty(), "isEmpty error")
	equalVal(t, false, pq.IsFull(), "isFull error")

	items := []*Item{
		{
			idx:      0,
			Priority: 5,
			Value:    "a",
		},
		{
			idx:      0,
			Priority: 3,
			Value:    "b",
		},
		{
			idx:      0,
			Priority: 8,
			Value:    "c",
		},
		{
			idx:      0,
			Priority: 6,
			Value:    "d",
		},
	}
	for _, item := range items {
		err := pq.Push(item)
		equalVal(t, nil, err, "push error")
	}
	equalVal(t, 5, pq.Len(), "len error")
	equalVal(t, pq.(*priorityQueue).queue.cap, pq.Cap(), "cap error")
	equalVal(t, false, pq.IsEmpty(), "isEmpty error")
	equalVal(t, false, pq.IsFull(), "isFull error")

	expect := []Item{
		{
			idx:      -1,
			Priority: 8,
			Value:    "c",
		},
		{
			idx:      -1,
			Priority: 7,
			Value:    "666",
		},
		{
			idx:      -1,
			Priority: 6,
			Value:    "d",
		},
		{
			idx:      -1,
			Priority: 5,
			Value:    "a",
		},
		{
			idx:      -1,
			Priority: 3,
			Value:    "b",
		},
	}

	//fmt.Println(pq.(*priorityQueue).queue.items)

	for _, item := range expect {
		val, err := pq.Pull()
		equalVal(t, nil, err, "pull error")
		equalVal(t, item, *(val.(*Item)), "pull error")
	}

	val, err := pq.Pull()
	equalVal(t, ErrIsEmpty, err, "pull error")
	equalVal(t, nil, val, "pull error")

}