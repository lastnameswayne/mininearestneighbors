package Set

type Set map[string]bool

func (s Set) GetRandom() string {
	for id, _ := range s {
		return id
	}

	return "-1"
}

func (s Set) UnsortedList() []string {
	res := []string{}
	for id, _ := range s {
		res = append(res, id)
	}

	return res
}

func (s Set) Add(element string) bool {
	s[element] = true
	return true
}
func (s Set) Delete(element string) bool {
	delete(s, element)
	return true
}

func (s Set) Has(element string) bool {
	_, ok := s[element]
	return ok
}
