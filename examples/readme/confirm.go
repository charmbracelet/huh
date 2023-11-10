package main

import (
	"github.com/charmbracelet/huh"
)

func main() {
	var happy bool

	confirm := huh.NewConfirm().
		Title("Are you sure?").
		Affirmative("Yes!").
		Negative("No.").
		Value(&happy)

	huh.NewForm(huh.NewGroup(confirm)).Run()
}
