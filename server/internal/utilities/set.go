package utilities

type void struct{} //empty structs occupy 0 memory

type Set[T comparable] map[T]void

func NewSet[T comparable](items []T, size ...int) Set[T] {
	var cap int
	if len(size) > 0 {
		cap = size[0]
	} else {
		cap = len(items)
	}
	s := make(Set[T], cap)
	for _, item := range items {
		s.Add(item)
	}
	return s
}

func (s Set[T]) Copy() Set[T] {
	newSet := NewSet([]T{})
	for elem := range s {
		newSet.Add(elem)
	}
	return newSet
}

func (s Set[T]) Has(v T) bool {
	_, ok := s[v]
	return ok
}

func (s Set[T]) Add(v T) {
	s[v] = void{}
}

func (s Set[T]) Remove(v T) {
	delete(s, v)
}

func (s Set[T]) Clear() {
	s = make(map[T]void)
}

func (s Set[T]) Size() int {
	return len(s)
}

func (s1 Set[T]) Is(s2 Set[T]) bool {
	if s1.Size() != s2.Size() {
		return false
	}
	for elem := range s1 {
		if !s2.Has(elem) {
			return false
		}
	}
	return true
}

func (s1 Set[T]) ToSlice() []T {
	slice := make([]T, 0, s1.Size())
	for elem := range s1 {
		slice = append(slice, elem)
	}
	return slice
}
