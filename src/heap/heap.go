package heap

import (
	"fmt"
)

type Heap []int

type heap interface {
	heapify()
	insert()
	delete()
}

func New() Heap {
	return Heap{}
}

func heapify(arr []int) Heap {
	h := New()
	for _, elem := range arr {
		h.insert(elem)
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
		if left < size && (*h)[left] < (*h)[smallest] {
			smallest = left
		}
		if right < size && (*h)[right] < (*h)[smallest] {
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
	for elem < parentVal {
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

func (h *Heap) insert(element int) Heap {
	*h = append(*h, element)
	h.up()
	return *h
}

func (h *Heap) delete(element int) int {
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

func (h *Heap) peek() int {
	if len(*h) < 1 {
		return 0
	}
	return (*h)[0]
}
