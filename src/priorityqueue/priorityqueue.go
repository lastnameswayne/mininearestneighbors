package priorityqueue

import "github.com/lastnameswayne/mininearestneighbors/src/heap"

type Element struct {
	weight  int
	element int
}

type priorityQueue struct {
	queue []Element
}

type PriorityQueue interface {
	Pop()
	Push()
}

func New(size int) *PriorityQueue {
	priorityQueue := priorityQueue{
		queue: heap.New(),
	}
	return &priorityQueue
}

func (p *PriorityQueue) Pop() (int, int) {
	h := heap.New()

}

func (p *PriorityQueue) Push(element int, weight int) {

}

func (p *PriorityQueue) Peak() {

}
