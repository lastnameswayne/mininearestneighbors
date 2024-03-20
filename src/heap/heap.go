package heap

import (
	g "github.com/lastnameswayne/mininearestneighbors/src/graph"
)

// this is a heap
// type is either min or max

type Heap struct {
	heap     []Element
	heapType Type
}

type Element struct {
	Weight float64
	Vertex g.Vertex
}

type Type int64

const (
	Min Type = 0
	Max Type = 1
)

func New(heapType Type) *Heap {
	return &Heap{
		heapType: heapType,
	}
}

// fills heap with arr
func (h *Heap) Heapify(arr []Element) *Heap {
	for _, elem := range arr {
		h.Insert(elem)
	}
	return h
}

func (h *Heap) down() {
	idx := 0
	for {
		left := 2*idx + 1
		right := 2*idx + 2

		smallest := idx
		if h.compare(left, smallest) {
			smallest = left
		}
		if h.compare(right, smallest) {
			smallest = right
		}

		if smallest == idx {
			break
		}

		h.swap(idx, smallest)

		idx = smallest
	}
}

func (h *Heap) compare(index, optimum int) bool {
	switch h.heapType {
	case Min:
		return index < len(h.heap) && h.heap[index].Weight < h.heap[optimum].Weight
	case Max:
		return index < len(h.heap) && h.heap[index].Weight > h.heap[optimum].Weight
	default:
		panic("unknown heap type ")
	}
}

func (h *Heap) up() {
	index := len(h.heap) - 1

	parent := (index - 1) / 2
	if parent < 0 {
		return
	}

	for h.compare(index, parent) {
		h.swap(index, int(parent))
		index = int(parent)

		parent = (index - 1) / 2
		if parent < 0 {
			break
		}
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

func (h *Heap) Sort(k int) []Element {
	res := []Element{}
	for i := 0; i < k; i++ {
		res = append(res, h.Delete())
	}

	hCopy := []Element{}
	copy(hCopy, h.heap)

	h.heap = hCopy

	return res
}

func (h *Heap) Elements() []Element {
	res := make([]Element, 0)
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

func (h *Heap) Size() int {
	if h == nil {
		return 0
	}
	return len(h.heap)
}

func (e *Element) GetWeight() float64 {
	if e == nil {
		return 100000000000000
	}

	return e.Weight
}
