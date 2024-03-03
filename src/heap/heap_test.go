package heap

import (
	"testing"

	g "github.com/lastnameswayne/mininearestneighbors/src/graph"
	"github.com/stretchr/testify/assert"
)

func TestHeapify(t *testing.T) {
	elems := []Element{
		{Weight: 1, Vertex: g.Vertex{Id: g.ID(1)}},
		{Weight: 5, Vertex: g.Vertex{Id: g.ID(2)}},
		{Weight: 3, Vertex: g.Vertex{Id: g.ID(3)}},
		{Weight: 6, Vertex: g.Vertex{Id: g.ID(4)}},
		{Weight: 7, Vertex: g.Vertex{Id: g.ID(5)}},
		{Weight: -1, Vertex: g.Vertex{Id: g.ID(6)}},
	}

	var heap = New(Min)

	heap.Heapify(elems)
	t.Log(heap)
	assert.Equal(t, -1.0, heap.heap[0].Weight)

	min := heap.Delete()
	t.Log(heap)

	assert.Equal(t, -1.0, min.Weight)
	assert.Equal(t, 5, len(heap.heap))

	heap = New(Min)
	elems = []Element{
		{Weight: 9, Vertex: g.Vertex{Id: g.ID(1)}},
		{Weight: 31, Vertex: g.Vertex{Id: g.ID(2)}},
		{Weight: 40, Vertex: g.Vertex{Id: g.ID(3)}},
		{Weight: 22, Vertex: g.Vertex{Id: g.ID(4)}},
		{Weight: 10, Vertex: g.Vertex{Id: g.ID(5)}},
		{Weight: 15, Vertex: g.Vertex{Id: g.ID(6)}},
		{Weight: 1, Vertex: g.Vertex{Id: g.ID(7)}},
		{Weight: 25, Vertex: g.Vertex{Id: g.ID(8)}},
		{Weight: 91, Vertex: g.Vertex{Id: g.ID(9)}},
	}
	heap.Heapify(elems)

	assert.Equal(t, 1.0, heap.heap[0].Weight)

}

func TestMaxHeap(t *testing.T) {
	elems := []Element{
		{Weight: 1, Vertex: g.Vertex{Id: g.ID(1)}},
		{Weight: 5, Vertex: g.Vertex{Id: g.ID(2)}},
		{Weight: 3, Vertex: g.Vertex{Id: g.ID(3)}},
		{Weight: 6, Vertex: g.Vertex{Id: g.ID(4)}},
		{Weight: 7, Vertex: g.Vertex{Id: g.ID(5)}},
		{Weight: -1, Vertex: g.Vertex{Id: g.ID(6)}},
	}

	var heap = New(Max)

	heap.Heapify(elems)
	t.Log(heap.heapType)
	assert.Equal(t, 7.0, heap.heap[0].Weight)

	heap = New(Max)
	elems = []Element{
		{Weight: 9, Vertex: g.Vertex{Id: g.ID(1)}},
		{Weight: 31, Vertex: g.Vertex{Id: g.ID(2)}},
		{Weight: 40, Vertex: g.Vertex{Id: g.ID(3)}},
		{Weight: 22, Vertex: g.Vertex{Id: g.ID(4)}},
		{Weight: 10, Vertex: g.Vertex{Id: g.ID(5)}},
		{Weight: 15, Vertex: g.Vertex{Id: g.ID(6)}},
		{Weight: 1, Vertex: g.Vertex{Id: g.ID(7)}},
		{Weight: 25, Vertex: g.Vertex{Id: g.ID(8)}},
		{Weight: 91, Vertex: g.Vertex{Id: g.ID(9)}},
	}
	heap.Heapify(elems)

	assert.Equal(t, 91.0, heap.heap[0].Weight)

	heap.Delete()
	assert.Equal(t, 40.0, heap.heap[0].Weight)
}
