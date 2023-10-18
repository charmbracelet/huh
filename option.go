package huh

// Option is a select option.
type Option[T any] struct {
	Key   string
	Value T
}

// NewOption returns a new select option.
func NewOption[T any](key string, value T) Option[T] {
	return Option[T]{Key: key, Value: value}
}
