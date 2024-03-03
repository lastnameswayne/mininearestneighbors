package heap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeapify(t *testing.T) {
	elems := []int{1, 5, 3, 6, 7, -1}

	var heap heap = elems

	heap = Heapify(elems)
	t.Log(heap)
	assert.Equal(t, -1, heap[0])

	min := heap.Delete(-1)
	t.Log(heap)

	assert.Equal(t, -1, min)
	assert.Equal(t, 5, len(heap))

	elems = []int{9, 31, 40, 22, 10, 15, 1, 25, 91}
	heap = Heapify(elems)

	assert.Equal(t, 1, heap[0])

}
