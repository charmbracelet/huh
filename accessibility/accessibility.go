package accessibility

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PromptInt prompts a user for an integer between a certain range.
//
// Given invalid input (non-integers, integers outside of the range), the user
// will continue to be reprompted until a valid input is given, ensuring that
// the return value is always valid.
func PromptInt(prompt string, min, max int) int {
	var (
		input  string
		choice int
	)

	validInt := func(s string) error {
		i, err := strconv.Atoi(s)
		if err != nil || i < min || i > max {
			return errors.New("invalid input. please try again")
		}
		return nil
	}

	input = PromptString(prompt, validInt)
	choice, _ = strconv.Atoi(input)
	return choice
}

// PromptBool prompts a user for a boolean value.
//
// Given invalid input (non-boolean), the user will continue to be reprompted
// until a valid input is given, ensuring that the return value is always valid.
func PromptBool() bool {
	validBool := func(s string) error {
		if len(s) == 1 && strings.Contains("yYnN", s) {
			return nil
		}
		return errors.New("invalid input. please try again")
	}

	input := PromptString("Choose [y/N]: ", validBool)
	return strings.ToLower(input) == "y"
}

// PromptString prompts a user for a string value and validates it against a
// validator function. It re-prompts the user until a valid input is given.
func PromptString(prompt string, validator func(input string) error) string {
	scanner := bufio.NewScanner(os.Stdin)

	var (
		valid bool
		input string
	)

	for !valid {
		fmt.Print(prompt)
		_ = scanner.Scan()
		input = scanner.Text()

		err := validator(input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		break
	}

	return input
}
