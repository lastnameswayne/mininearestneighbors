package Set

type Set map[int]bool

func (s Set) GetRandom() int {
	for id, _ := range s {
		return id
	}

	return -1
}

func (s Set) UnsortedList() []int {
	res := []int{}
	for id, _ := range s {
		res = append(res, id)
	}

	return res
}

func (s Set) Add(element int) bool {
	s[element] = true
	return true
}
func (s Set) Delete(element int) bool {
	delete(s, element)
	return true
}

func (s Set) Has(element int) bool {
	_, ok := s[element]
	return ok
}
