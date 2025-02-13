package huh

import "fmt"

// Option is an option for select fields.
type Option[T comparable] struct {
	Key      string
	Value    T
	selected bool
}

// NewOptions returns new options from a list of values.
func NewOptions[T comparable](values ...T) []Option[T] {
	options := make([]Option[T], len(values))
	for i, o := range values {
		options[i] = Option[T]{
			Key:   fmt.Sprint(o),
			Value: o,
		}
	}
	return options
}

// NewOption returns a new select option.
func NewOption[T comparable](key string, value T) Option[T] {
	return Option[T]{Key: key, Value: value}
}

// Selected sets whether the option is currently selected.
func (o Option[T]) Selected(selected bool) Option[T] {
	o.selected = selected
	return o
}

// String returns the key of the option.
func (o Option[T]) String() string {
	return o.Key
}

// TableOption is an option for table fields.
// The T type represents a row as a struct.
// The K type is the key used to identify the row.
type TableOption[T any, K comparable] struct {
	key      func(T) K
	Value    T
	selected bool
}

// NewTableOptions returns new options from a list of values.
func NewTableOptions[T any, K comparable](key func(T) K, values ...T) []TableOption[T, K] {
	options := make([]TableOption[T, K], len(values))
	for i, o := range values {
		options[i] = TableOption[T, K]{
			key:   key,
			Value: o,
		}
	}
	return options
}

// NewTableOption returns a new table option.
func NewTableOption[T any, K comparable](key func(T) K, value T) TableOption[T, K] {
	return TableOption[T, K]{
		key:   key,
		Value: value,
	}
}

// Key returns the key value for this option.
func (o TableOption[T, K]) Key() K {
	return o.key(o.Value)
}

// Selected sets whether the option is currently selected.
func (o TableOption[T, K]) Selected(selected bool) TableOption[T, K] {
	o.selected = selected
	return o
}
