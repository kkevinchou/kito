package priorityqueue

import "container/heap"

// An Item is something we manage in a priority queue.
type Item struct {
	value    interface{} // The value of the item; arbitrary.
	priority float64     // The priority of the item in the queue.
}

// A Heap implements heap.Interface and holds Items.
type Heap []*Item

func (h Heap) Len() int { return len(h) }

func (h *Heap) Less(i, j int) bool {
	return (*h)[i].priority < (*h)[j].priority
}

func (h *Heap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *Heap) Push(x interface{}) {
	item := x.(*Item)
	*h = append(*h, item)
}

func (h *Heap) Pop() interface{} {
	n := len(*h)
	item := (*h)[n-1]
	*h = (*h)[0 : n-1]
	return item
}

type PriorityQueue struct {
	h *Heap
}

func New() *PriorityQueue {
	pq := PriorityQueue{h: &Heap{}}
	heap.Init(pq.h)
	return &pq
}

func (pq *PriorityQueue) Push(x interface{}, priority float64) {
	heap.Push(pq.h, &Item{value: x, priority: priority})
}

func (pq *PriorityQueue) Pop() interface{} {
	item := heap.Pop(pq.h).(*Item)
	return item.value
}

func (pq *PriorityQueue) Len() int {
	return pq.h.Len()
}

func (pq *PriorityQueue) Empty() bool {
	return pq.Len() == 0
}
