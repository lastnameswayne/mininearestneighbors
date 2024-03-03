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

func New(query g.Vertex, sort heap.Type) *PriorityQueue {
	PriorityQueue := PriorityQueue{
		queue: *heap.New(sort),
		query: query,
	}
	return &PriorityQueue
}

func (p *PriorityQueue) Pop() heap.Element {
	return p.queue.Delete()
}

func (p *PriorityQueue) Push(vertex g.Vertex) {
	weight := v.Distance(vertex.Vector, p.query.Vector)

	elem := heap.Element{
		Weight: weight,
		Vertex: vertex,
	}

	p.queue.Insert(elem)
}

func (p *PriorityQueue) Peek() *heap.Element {
	return p.queue.Peek()

}

func (p *PriorityQueue) Size() int {
	return p.queue.Size()
}

func (p *PriorityQueue) Elements() []heap.Element {
	return p.queue.Elements()
}
