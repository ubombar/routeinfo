package structures

// Set is a generic set data structure.
type Set[T comparable] struct {
	elements map[T]struct{}
}

// NewSet creates a new Set.
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{elements: make(map[T]struct{})}
}

// NewSetWithSize creates a new Set with specified size.
func NewSetWithSize[T comparable](size int) *Set[T] {
	return &Set[T]{elements: make(map[T]struct{}, size)}
}

// Add inserts an element into the set.
func (s *Set[T]) Add(value T) {
	s.elements[value] = struct{}{}
}

// Remove deletes an element from the set.
func (s *Set[T]) Remove(value T) {
	delete(s.elements, value)
}

// Contains checks if an element is in the set.
func (s *Set[T]) Contains(value T) bool {
	_, exists := s.elements[value]
	return exists
}

// Elements returns a slice of all elements in the set.
func (s *Set[T]) Elements() []T {
	result := make([]T, 0, len(s.elements))
	for key := range s.elements {
		result = append(result, key)
	}
	return result
}

// Size returns the number of elements in the set.
func (s *Set[T]) Size() int {
	return len(s.elements)
}
