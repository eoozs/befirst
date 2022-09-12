package set

type Set[T comparable] struct {
	values map[T]struct{}
}

func New[T comparable]() *Set[T] {
	return &Set[T]{
		values: make(map[T]struct{}),
	}
}

func (s *Set[T]) Add(values ...T) {
	for _, val := range values {
		s.values[val] = struct{}{}
	}
}

func (s *Set[T]) Clear() {
	s.values = make(map[T]struct{})
}

func (s *Set[T]) Remove(val T) {
	delete(s.values, val)
}

func (s *Set[T]) Contains(val T) bool {
	_, contains := s.values[val]
	return contains
}

func (s *Set[T]) Size() int {
	return len(s.values)
}

func (s *Set[T]) ToSlice() []T {
	list := make([]T, 0, s.Size())
	for val := range s.values {
		list = append(list, val)
	}
	return list
}
