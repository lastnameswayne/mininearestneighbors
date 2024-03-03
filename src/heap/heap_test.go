package heap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeapify(t *testing.T) {
	elems := []Element{
		{Weight: 1, Element: 1}, {Weight: 5, Element: 2}, {Weight: 3, Element: 3}, {Weight: 6, Element: 4}, {Weight: 7, Element: 5}, {Weight: -1, Element: 6},
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
	elems = []Element{{Weight: 9, Element: 1}, {Weight: 31, Element: 2}, {Weight: 40, Element: 3}, {Weight: 22, Element: 4}, {Weight: 10, Element: 5}, {Weight: 15, Element: 6}, {Weight: 1, Element: 7}, {Weight: 25, Element: 8}, {Weight: 91, Element: 9}}

	heap.Heapify(elems)

	assert.Equal(t, 1.0, heap.heap[0].Weight)

}

func TestMaxHeap(t *testing.T) {
	elems := []Element{
		{Weight: 1, Element: 1}, {Weight: 5, Element: 2}, {Weight: 3, Element: 3}, {Weight: 6, Element: 4}, {Weight: 7, Element: 5}, {Weight: -1, Element: 6},
	}

	var heap = New(Max)

	heap.Heapify(elems)
	t.Log(heap.heapType)
	assert.Equal(t, 7.0, heap.heap[0].Weight)

	heap = New(Max)
	elems = []Element{{Weight: 9, Element: 1}, {Weight: 31, Element: 2}, {Weight: 40, Element: 3}, {Weight: 22, Element: 4}, {Weight: 10, Element: 5}, {Weight: 15, Element: 6}, {Weight: 1, Element: 7}, {Weight: 25, Element: 8}, {Weight: 91, Element: 9}}

	heap.Heapify(elems)

	assert.Equal(t, 91.0, heap.heap[0].Weight)

	heap.Delete()
	assert.Equal(t, 40.0, heap.heap[0].Weight)
}
