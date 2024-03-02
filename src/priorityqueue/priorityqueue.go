package priorityqueue

type Element struct {
	weight  int
	element int
}

type PriorityQueue struct {
	queue []Element
	size  int
}

type queue interface {
	Pop()
	Push()
}

type binaryTree struct {
	nodes []node
}

type node struct {
	left  *node
	right *node
}

func New(size int) *PriorityQueue {
	priorityQueue := PriorityQueue{
		queue: []Element{},
		size:  size,
	}
	return &priorityQueue
}

func (p *PriorityQueue) Pop() (int, int) {

}

func (p *PriorityQueue) Push(element int, weight int) {

}