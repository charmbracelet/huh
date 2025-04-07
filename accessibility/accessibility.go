// Package accessibility provides accessible functions to capture user input.
package accessibility

import (
	"bufio"
	"cmp"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

// PromptInt prompts a user for an integer between a certain range.
//
// Given invalid input (non-integers, integers outside of the range), the user
// will continue to be reprompted until a valid input is given, ensuring that
// the return value is always valid.
func PromptInt(prompt string, low, high int, opts ...Option) int {
	var (
		input  string
		choice int
	)

	validInt := func(s string) error {
		i, err := strconv.Atoi(s)
		if err != nil || i < low || i > high {
			return errors.New("invalid input. please try again")
		}
		return nil
	}

	input = PromptString(prompt, validInt, opts...)
	choice, _ = strconv.Atoi(input)
	return choice
}

func parseBool(s string, defaultValue bool) (bool, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return defaultValue, nil
	}

	if slices.Contains([]string{"y", "yes"}, s) {
		return true, nil
	}

	// As a special case, we default to "" to no since the usage of this
	// function suggests N is the default.
	if slices.Contains([]string{"n", "no"}, s) {
		return false, nil
	}

	return false, errors.New("invalid input. please try again")
}

// PromptBool prompts a user for a boolean value.
//
// Given invalid input (non-boolean), the user will continue to be reprompted
// until a valid input is given, ensuring that the return value is always valid.
func PromptBool(opts ...Option) bool {
	options := eval(opts)
	defaultValue := options.defaultValue.(bool)
	validBool := func(s string) error {
		_, err := parseBool(s, defaultValue)
		return err
	}

	chooseStr := "y/N"
	if defaultValue {
		chooseStr = "Y/n"
	}
	input := PromptString("Choose ["+chooseStr+"]: ", validBool, opts...)
	b, _ := parseBool(input, defaultValue)
	return b
}

// PromptString prompts a user for a string value and validates it against a
// validator function. It re-prompts the user until a valid input is given.
func PromptString(prompt string, validator func(input string) error, opts ...Option) string {
	w, r := ioFor(eval(opts))
	scanner := bufio.NewScanner(r)

	var (
		valid bool
		input string
	)

	for !valid {
		_, _ = fmt.Fprint(w, prompt)
		if !scanner.Scan() {
			// no way to bubble up errors or signal cancellation
			// but the program is probably not continuing if
			// stdin sent EOF
			break
		}
		input = scanner.Text()

		err := validator(input)
		if err != nil {
			_, _ = fmt.Fprintln(w, err)
			continue
		}

		break
	}

	return input
}

type options struct {
	w            io.Writer
	r            io.Reader
	defaultValue any
}

func eval(opts []Option) options {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func ioFor(o options) (io.Writer, io.Reader) {
	return cmp.Or[io.Writer](o.w, os.Stdout),
		cmp.Or[io.Reader](o.r, os.Stdin)
}

// Option sets the options for the accessibility operations.
type Option func(*options)

// Output sets the output writer for the accessibility operations.
func Output(w io.Writer) Option {
	return func(o *options) {
		o.w = w
	}
}

// Input sets the input writer for the accessibility operations.
func Input(r io.Reader) Option {
	return func(o *options) {
		o.r = r
	}
}

// DefaultValue sets the default value of the field.
func DefaultValue(defaultValue any) Option {
	return func(o *options) {
		o.defaultValue = defaultValue
	}
}
