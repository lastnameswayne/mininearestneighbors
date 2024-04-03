package Set

type T any

type Set[T comparable] map[T]bool

func New[T comparable]() Set[T] {
	return make(Set[T])
}

func (s Set[T]) GetRandom() T {
	for id, _ := range s {
		return id
	}

	var emptyReturn T
	return emptyReturn
}

func (s Set[T]) UnsortedList() []T {
	res := []T{}
	for id, _ := range s {
		res = append(res, id)
	}

	return res
}

func (s Set[T]) Add(element T) bool {
	s[element] = true
	return true
}
func (s Set[T]) Delete(element T) bool {
	delete(s, element)
	return true
}

func (s Set[T]) Has(element T) bool {
	_, ok := s[element]
	return ok
}
