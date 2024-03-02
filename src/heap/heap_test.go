package heap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeapify(t *testing.T) {
	elems := []int{1, 5, 3, 6, 7, -1}

	var heap Heap = elems

	heap = heapify(elems)
	t.Log(heap)
	assert.Equal(t, -1, heap[0])

	min := heap.delete(-1)
	t.Log(heap)

	assert.Equal(t, -1, min)

	elems = []int{9, 31, 40, 22, 10, 15, 1, 25, 91}
	heap = heapify(elems)

	assert.Equal(t, 1, heap[0])

}
