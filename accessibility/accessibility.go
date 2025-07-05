// Package accessibility provides accessible functions to capture user input.
//
// Deprecated: use [internal/accessibility] instead.
package accessibility

import (
	"os"

	"github.com/charmbracelet/huh/internal/accessibility"
)

// PromptInt prompts a user for an integer between a certain range.
//
// Given invalid input (non-integers, integers outside of the range), the user
// will continue to be reprompted until a valid input is given, ensuring that
// the return value is always valid.
//
// Deprecated: use [accessibility.PromptInt] instead.
func PromptInt(prompt string, low, high int) int {
	return accessibility.PromptInt(os.Stdout, os.Stdin, prompt, low, high, nil)
}

// PromptBool prompts a user for a boolean value.
//
// Given invalid input (non-boolean), the user will continue to be reprompted
// until a valid input is given, ensuring that the return value is always valid.
//
// Deprecated: use [accessibility.PromptBool] instead.
func PromptBool() bool {
	return accessibility.PromptBool(os.Stdout, os.Stdin, "Choose [y/N]: ", false)
}

// PromptString prompts a user for a string value and validates it against a
// validator function. It re-prompts the user until a valid input is given.
//
// Deprecated: use [accessibility.PromptString] instead.
func PromptString(prompt string, validator func(input string) error) string {
	return accessibility.PromptString(os.Stdout, os.Stdin, prompt, "", validator)
}
