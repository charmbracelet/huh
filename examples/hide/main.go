package main

import (
	"fmt"

	"github.com/charmbracelet/huh/v2"
)

func main() {
	var isAllergic bool
	var allergies string

	huh.NewForm(
		huh.NewGroup(huh.NewNote().Title("Just for fun!")).WithHideFunc(func() bool { return true }),
		huh.NewGroup(huh.NewNote().Title("Just for fun!")).WithHide(true),

		huh.NewGroup(huh.NewConfirm().
			Title("Do you have any allergies?").
			Description("If so, please list them.").
			Value(&isAllergic)),
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
