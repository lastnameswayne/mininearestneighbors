package Set

type IntSet interface {
	Add(element int) bool
	Delete(element int) bool
	Has(element int) bool
}

type Set map[int]bool

func NewSet() Set {
	return map[int]bool{}
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
