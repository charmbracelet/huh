package huh

import (
	"fmt"
	"unicode/utf8"
)

// ErrInputEmpty is an error indicating that the input cannot be empty.
var ErrInputEmpty = fmt.Errorf("input cannot be empty")

// ValidateNotEmpty checks if the input is not empty.
func ValidateNotEmpty() func(s string) error {
	return func(s string) error {
		if s == "" {
			return ErrInputEmpty
		}
		return nil
	}
}

// ValidateLength checks if the length of the input is within the specified range.
func ValidateLength(min, max int) func(s string) error {
	return func(s string) error {
		length := utf8.RuneCountInString(s)
		if length < min || length > max {
			return fmt.Errorf("input must be between %d and %d", min, max)
		}
		return nil
	}
}

// ValidateOneOf checks if a string is one of the specified options.
func ValidateOneOf(options ...string) func(string) error {
	validOptions := make(map[string]struct{})
	for _, option := range options {
		validOptions[option] = struct{}{}
	}

	return func(value string) error {
		if _, ok := validOptions[value]; !ok {
			return fmt.Errorf("invalid option: %s", value)
		}
		return nil
	}
}
