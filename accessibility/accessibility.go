package accessibility

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func PromptInt(prompt string, min, max int) int {
	var (
		input  string
		valid  bool
		choice int
		err    error
	)

	for !valid {
		input = PromptString(prompt)
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
		input = PromptString("Choose [y/N]: ")

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

func PromptString(prompt string) string {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print(prompt)
	_ = scanner.Scan()

	text := scanner.Text()
	return text
}
