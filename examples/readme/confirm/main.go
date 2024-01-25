package main

import (
	"github.com/charmbracelet/huh"
)

func main() {
	var happy bool

	confirm := huh.NewConfirm().
		Title("Are you sure? ").
		Description("Please confirm. ").
		Affirmative("Yes!").
		Negative("No.").
		Inline(true).
		Value(&happy)

	huh.NewForm(huh.NewGroup(confirm)).Run()
}
