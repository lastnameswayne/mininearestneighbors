package heap

import (
	"fmt"
	"math"
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

func (h *Heap) heapify() {
	elem := (*h)[i]
	parent := math.Floor(float64((i - 1) / 2))
	if parent < 0 {
		continue
	}

	newIndex := i
	parentVal := (*h)[int(parent)]
	for elem < parentVal {
		fmt.Println(elem, parentVal)

		h.swap(newIndex, int(parent))
		newIndex = int(parent)

		parent = math.Floor(float64((newIndex - 1) / 2))
		if parent < 0 {
			break
		}
		parentVal = (*h)[int(parent)]
	}
}

func (h Heap) insert(element int) Heap {

	h = append(h, element)
	return h
}

func (h Heap) delete(element int) {
}

func (h *Heap) swap(idx1, idx2 int) {
	(*h)[idx1], (*h)[idx2] = (*h)[idx2], (*h)[idx1]
}
