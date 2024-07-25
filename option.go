package huh

import (
	"fmt"
)

// Option is an option for select fields.
type Option[T comparable] struct {
	Value T
}

// NewOptions returns new options from a list of values.
func NewOptions[T comparable](values ...T) []Option[T] {
	options := make([]Option[T], len(values))
	for i, o := range values {
		options[i] = Option[T]{
			Value: o,
		}
	}
	return options
}

// String returns the string representation of the Option.
func (o Option[T]) String() string {
	return fmt.Sprint(o.Value)
}

// NewOption returns a new select option.
func NewOption[T comparable](value T) Option[T] {
	return Option[T]{Value: value}
}
