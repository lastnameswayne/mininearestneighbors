package priorityqueue

import (
	g "github.com/lastnameswayne/mininearestneighbors/src/graph"
	heap "github.com/lastnameswayne/mininearestneighbors/src/heap"
	s "github.com/lastnameswayne/mininearestneighbors/src/set"
	v "github.com/lastnameswayne/mininearestneighbors/src/vector"
)

// the weight is automatically determined
// by comparing to the query
// you define the query on initialization
type PriorityQueue struct {
	queue heap.Heap
	query g.Vertex
	set   s.Set
}

func New(query g.Vertex, sort heap.Type) *PriorityQueue {
	PriorityQueue := PriorityQueue{
		queue: *heap.New(sort),
		query: query,
		set:   s.Set{},
	}
	return &PriorityQueue
}

func (p *PriorityQueue) Pop() heap.Element {
	deleted := p.queue.Delete()
	p.set.Delete(int(deleted.Vertex.Id))
	return deleted
}

func (p *PriorityQueue) Push(vertex g.Vertex) {
	if p.set.Has(int(vertex.Id)) {
		return
	}
	weight := v.Distance(vertex.Vector, p.query.Vector)

	elem := heap.Element{
		Weight: weight,
		Vertex: vertex,
	}

	p.queue.Insert(elem)
	p.set.Add(int(vertex.Id))
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
