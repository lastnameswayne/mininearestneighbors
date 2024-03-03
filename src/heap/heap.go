package heap

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
	if h.heapType == Min {
		return index < len(h.heap) && h.heap[index].Weight < h.heap[optimum].Weight
	}
	return false
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
