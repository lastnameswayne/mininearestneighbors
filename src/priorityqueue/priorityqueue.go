package priorityqueue

import "github.com/lastnameswayne/mininearestneighbors/src/heap"

type priorityQueue struct {
	queue heap.Heap
}

type PriorityQueue interface {
	Pop()
	Push()
}

func New() *priorityQueue {
	priorityQueue := priorityQueue{
		queue: heap.New(),
	}
	return &priorityQueue
}

func (p *priorityQueue) Pop() (int, int) {
	h := heap.New()
	h.Delete()

}

func (p *priorityQueue) Push(element int, weight int) {

}

func (p *priorityQueue) Peak() {

}
