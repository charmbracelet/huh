package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

func main() {
	var isAllergic bool
	var allergies string

	huh.NewForm(
		huh.NewGroup(huh.NewNote().Title("Just for fun!")).WithHideFunc(func() bool { return true }),
		huh.NewGroup(huh.NewNote().Title("Just for fun!")).WithHide(true),

		huh.NewGroup(huh.NewConfirm().Title("Are you allergic to anything?").Value(&isAllergic)),
		huh.NewGroup(
			huh.NewText().
				Title("Allergies").
				Description("Please list all your allergies...").
				Value(&allergies),
		).WithHideFunc(func() bool {
			return !isAllergic
		}),
		huh.NewGroup(huh.NewNote().Title("Invisible")).WithHide(true),
	).Run()

	if isAllergic {
		fmt.Println(allergies)
	}
}
