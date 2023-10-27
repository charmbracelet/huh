package huh

// Option is an option for select fields.
type Option[T any] struct {
	Key      string
	Value    T
	selected bool
}

// NewOption returns a new select option.
func NewOption[T any](key string, value T) Option[T] {
	return Option[T]{Key: key, Value: value}
}

// Selected sets whether the option is currently selected.
func (o Option[T]) Selected(selected bool) Option[T] {
	o.selected = selected
	return o
}
