package heap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeapify(t *testing.T) {
	elems := []Element{
		{weight: 1, element: 1}, {weight: 5, element: 2}, {weight: 3, element: 3}, {weight: 6, element: 4}, {weight: 7, element: 5}, {weight: -1, element: 6},
	}

	var heap Heap = elems

	heap = Heapify(elems)
	t.Log(heap)
	assert.Equal(t, -1, heap[0].weight)

	min := heap.Delete()
	t.Log(heap)

	assert.Equal(t, -1, min.weight)
	assert.Equal(t, 5, len(heap))

	elems = []Element{{weight: 9, element: 1}, {weight: 31, element: 2}, {weight: 40, element: 3}, {weight: 22, element: 4}, {weight: 10, element: 5}, {weight: 15, element: 6}, {weight: 1, element: 7}, {weight: 25, element: 8}, {weight: 91, element: 9}}

	heap = Heapify(elems)

	assert.Equal(t, 1, heap[0].weight)

}
