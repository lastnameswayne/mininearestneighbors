package heap

import (
	"fmt"
)

// this is a minheap

type Heap []Element

type Element struct {
	Weight  float64
	Element int
}

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
	size := len(*h)
	for {
		left := 2*idx + 1
		right := 2*idx + 2

		smallest := idx
		if left < size && (*h)[left].Weight < (*h)[smallest].Weight {
			smallest = left
		}
		if right < size && (*h)[right].Weight < (*h)[smallest].Weight {
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
	i := len(*h) - 1

	elem := (*h)[i]
	parent := (i - 1) / 2
	if parent < 0 {
		return
	}

	newIndex := i
	parentVal := (*h)[int(parent)]
	for elem.Weight < parentVal.Weight {
		fmt.Println(elem, parentVal)

		h.swap(newIndex, int(parent))
		newIndex = int(parent)

		parent = (newIndex - 1) / 2
		if parent < 0 {
			break
		}
		parentVal = (*h)[int(parent)]
	}
}

func (h *Heap) Insert(element Element) Heap {
	*h = append(*h, element)
	h.up()
	return *h
}

func (h *Heap) Delete() Element {
	min := (*h)[0]

	//push root to back and cut it off
	lastElementIndex := len(*h) - 1
	val := (*h)[lastElementIndex]
	(*h)[0] = val
	*h = (*h)[:lastElementIndex]

	h.down()

	return min
}

func (h *Heap) swap(idx1, idx2 int) {
	(*h)[idx1], (*h)[idx2] = (*h)[idx2], (*h)[idx1]
}

func (h *Heap) Peek() *Element {
	if len(*h) < 1 {
		return nil
	}
	return &(*h)[0]
}
