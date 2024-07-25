package selector

// Selector is a helper type for selecting items.
type Selector struct {
	index          int
	numberOfFields int
}

// NewSelector creates a new item selector.
func NewSelector(n int) *Selector {
	return &Selector{
		numberOfFields: n,
	}
}

// Next moves the selector to the next item.
func (s *Selector) Next() {
	if s.index < s.numberOfFields-1 {
		s.index++
	}
}

// Prev moves the selector to the previous item.
func (s *Selector) Prev() {
	if s.index > 0 {
		s.index--
	}
}

// OnFirst returns true if the selector is on the first item.
func (s *Selector) OnFirst() bool {
	return s.index == 0
}

// OnLast returns true if the selector is on the last item.
func (s *Selector) OnLast() bool {
	return s.index == s.numberOfFields-1
}

// Selected returns the index of the current selected item.
func (s *Selector) Selected() int {
	return s.index
}

// Totoal returns the total number of items.
func (s *Selector) Total() int {
	return s.numberOfFields
}

// SetSelected sets the selected item.
func (s *Selector) SetSelected(i int) {
	if i < 0 || i >= s.numberOfFields {
		return
	}
	s.index = i
}
