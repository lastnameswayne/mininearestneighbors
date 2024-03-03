package priorityqueue

import (
	g "github.com/lastnameswayne/mininearestneighbors/src/graph"
	heap "github.com/lastnameswayne/mininearestneighbors/src/heap"
	v "github.com/lastnameswayne/mininearestneighbors/src/vector"
)

// the weight is automatically determined
// by comparing to the query
// you define the query on initialization
type PriorityQueue struct {
	queue heap.Heap
	query g.Vertex
}

func New(query g.Vertex) *PriorityQueue {
	PriorityQueue := PriorityQueue{
		queue: heap.New(),
		query: query,
	}
	return &PriorityQueue
}

func (p *PriorityQueue) Pop() heap.Element {
	h := heap.New()
	return h.Delete()
}

func (p *PriorityQueue) Push(vector g.Vertex) {
	weight := v.Distance(vector.Vector, p.query.Vector)

	elem := heap.Element{
		Weight:  weight,
		Element: int(vector.Id),
	}

	p.queue.Insert(elem)
}

func (p *PriorityQueue) Peek() *heap.Element {
	return p.queue.Peek()

}

func (p *PriorityQueue) Size() int {
	return len(p.queue)
}
