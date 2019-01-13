package containers

import "bytes"

// Set struct provides the definition for a string set.
type Set struct {
	elements map[string]struct{}
}

// NewSet creates and returns an empty set.
func NewSet() (s *Set) {
	s = &Set{make(map[string]struct{})}
	return
}

// Add adds an element to the set.
func (s *Set) Add(elem string) {
	s.elements[elem] = struct{}{}
}

// Delete removes an element from the set.
func (s *Set) Delete(elem string) {
	delete(s.elements, elem)
}

// AddAll transforms s into the union of s and s2.
func (s *Set) AddAll(s2 *Set) {
	for k := range s2.elements {
		s.elements[k] = struct{}{}
	}
}

// RetainAll transforms s into the intersection of s and s2.
func (s *Set) RetainAll(s2 *Set) {
	for k := range s.elements {
		if !s2.Contains(k) {
			delete(s.elements, k)
		}
	}
}

// DeleteAll transforms s into the difference between s and s2.
func (s *Set) DeleteAll(s2 *Set) {
	for k := range s.elements {
		if s2.Contains(k) {
			delete(s.elements, k)
		}
	}
}

// Contains returns true if the element
// is present in the set, false otherwise.
func (s *Set) Contains(elem string) (present bool) {
	_, present = s.elements[elem]
	return
}

// Iterator returns a iterator function for the set.
// WARNING: Not safe for inserts or deletions of elements during the iteration.
func (s *Set) Iterator() func() (elem string, ok bool) {
	keys := make([]string, len(s.elements))
	i := 0
	for k := range s.elements {
		keys[i] = k
		i++
	}

	i = -1
	return func() (elem string, ok bool) {
		i++
		if i == len(keys) {
			ok = false
			return
		}
		ok = true
		elem = keys[i]
		return
	}
}

// Len returns the lenght of a set.
func (s *Set) Len() (l int) {
	l = len(s.elements)
	return
}

func (s *Set) String() string {
	var b bytes.Buffer

	b.WriteString("{ ")
	for k := range s.elements {
		b.WriteString(k)
		b.WriteString(" ")
	}
	b.WriteString("}")

	return b.String()
}
