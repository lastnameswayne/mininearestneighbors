package heap

import (
	"fmt"
)

// this is a heap
// type is either min or max

type Heap struct {
	heap     []Element
	heapType Type
}

type Element struct {
	Weight  float64
	Element int
}

type Type int64

const (
	Min Type = 0
	Max Type = 1
)

func New() Heap {
	return Heap{}
}

func Heapify(arr []Element) Heap {
	h := New()
	for _, elem := range arr {
		h.Insert(elem)
	}
	return h
}

func (h *Heap) down() {
	idx := 0
	size := len(h.heap)
	for {
		left := 2*idx + 1
		right := 2*idx + 2

		smallest := idx
		if left < size && h.heap[left].Weight < h.heap[smallest].Weight {
			smallest = left
		}
		if right < size && h.heap[right].Weight < h.heap[smallest].Weight {
			smallest = right
		}

		if smallest == idx {
			break
		}

		h.swap(idx, smallest)

		idx = smallest
	}
}

func (h *Heap) up() {
	i := len(h.heap) - 1

	elem := h.heap[i]
	parent := (i - 1) / 2
	if parent < 0 {
		return
	}

	newIndex := i
	parentVal := h.heap[int(parent)]
	for elem.Weight < parentVal.Weight {
		fmt.Println(elem, parentVal)

		h.swap(newIndex, int(parent))
		newIndex = int(parent)

		parent = (newIndex - 1) / 2
		if parent < 0 {
			break
		}
		parentVal = h.heap[int(parent)]
	}
}

func (h *Heap) Insert(element Element) Heap {
	h.heap = append(h.heap, element)
	h.up()
	return *h
}

func (h *Heap) Delete() Element {
	min := h.heap[0]

	//push root to back and cut it off
	lastElementIndex := len(h.heap) - 1
	val := h.heap[lastElementIndex]
	h.heap[0] = val
	h.heap = h.heap[:lastElementIndex]

	h.down()

	return min
}

func (h *Heap) Elements() []Element {
	res := make([]Element, len(h.heap))
	res = append(res, h.heap...)
	return res
}

func (h *Heap) swap(idx1, idx2 int) {
	h.heap[idx1], h.heap[idx2] = h.heap[idx2], h.heap[idx1]
}

func (h *Heap) Peek() *Element {
	if len(h.heap) < 1 {
		return nil
	}
	return &h.heap[0]
}
