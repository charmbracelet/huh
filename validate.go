package huh

import (
	"fmt"
)

// ValidateNotEmpty checks if the input is not empty.
func ValidateNotEmpty(s string) error {
	if s == "" {
		return fmt.Errorf("input cannot be empty")
	}
	return nil
}

// ValidateLength checks if the input has a length within the specified range.
func ValidateLength(min, max int) func([]string) error {
	return func(input []string) error {
		length := len(input)
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
			return fmt.Errorf("invalid option")
		}
		return nil
	}
}