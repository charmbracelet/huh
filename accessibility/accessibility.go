// Package accessibility provides accessible functions to capture user input.
package accessibility

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

// PromptInt prompts a user for an integer between a certain range.
//
// Given invalid input (non-integers, integers outside of the range), the user
// will continue to be reprompted until a valid input is given, ensuring that
// the return value is always valid.
func PromptInt(w io.Writer, r io.Reader, prompt string, low, high int) int {
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

	input = PromptString(w, r, prompt, validInt)
	choice, _ = strconv.Atoi(input)
	return choice
}

func parseBool(s string, val bool) (bool, error) {
	s = strings.ToLower(s)

	if s == "" {
		return val, nil
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
func PromptBool(w io.Writer, r io.Reader, val bool) bool {
	validBool := func(s string) error {
		_, err := parseBool(s, val)
		return err
	}

	options := "y/N"
	if val {
		options = "Y/n"
	}
	input := PromptString(w, r, "Choose ["+options+"]: ", validBool)
	b, _ := parseBool(input, val)
	return b
}

// PromptString prompts a user for a string value and validates it against a
// validator function. It re-prompts the user until a valid input is given.
func PromptString(w io.Writer, r io.Reader, prompt string, validator func(input string) error) string {
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
