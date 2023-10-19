package accessibility

import (
	"fmt"
	"strconv"
)

func PromptInt(min, max int) int {
	var (
		input  string
		valid  bool
		choice int
		err    error
	)

	for !valid {
		fmt.Print("Choose: ")

		// We scan the entire line so that if the input is invalid, we return only
		// one error message instead of one for each scan.
		fmt.Scanln(&input)

		choice, err = strconv.Atoi(input)

		if err != nil {
			fmt.Println("Invalid input. Please try again.")
			continue
		}

		if choice < min || choice > max {
			fmt.Println("Invalid input. Please try again.")
			continue
		}

		break
	}

	return choice
}

func PromptBool() bool {
	var (
		input  string
		valid  bool
		choice bool
	)

	for !valid {
		fmt.Print("Choose [y/N]: ")

		// We scan the entire line so that if the input is invalid, we return only
		// one error message instead of one for each scan.
		fmt.Scanln(&input)

		if input == "y" || input == "Y" {
			choice = true
		} else if input == "n" || input == "N" {
			choice = false
		} else {
			fmt.Println("Invalid input. Please try again.")
			continue
		}

		break
	}

	return choice
}

func PromptString() string {
	var (
		input string
	)

	fmt.Print("> ")
	fmt.Scanln(&input)

	return input
}
