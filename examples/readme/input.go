package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

func isFood(_ string) error {
	return nil
}

func main() {
	var lunch string
	huh.NewInput().
		Title("What's for lunch?").
		Prompt("? ").
		Validate(isFood).
		Value(&lunch).
		Run()

	fmt.Printf("Yummy, %s!\n", lunch)
}
