package heap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeapify(t *testing.T) {
	elems := []int{1, 5, 3, 6, 7, -1}

	var heap Heap = elems

	heap.heapify()
	t.Log(heap)
	assert.Equal(t, -1, heap[0])

	elems = []int{12, 11, 13, 5, 6, 7}
	heap = elems
	heap.heapify()

	expected := []int{5, 6, 7, 11, 12, 13}
	var expectedHeap Heap = expected
	assert.Equal(t, expectedHeap, heap)

}
