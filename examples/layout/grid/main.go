package main

import "github.com/charmbracelet/huh/v2"

func main() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("First"),
			huh.NewInput().Title("Second"),
			huh.NewInput().Title("Third"),
		),
		huh.NewGroup(
			huh.NewInput().Title("Fourth"),
			huh.NewInput().Title("Fifth"),
			huh.NewInput().Title("Sixth"),
		),
		huh.NewGroup(
			huh.NewInput().Title("Seventh"),
			huh.NewInput().Title("Eigth"),
			huh.NewInput().Title("Nineth"),
			huh.NewInput().Title("Tenth"),
		),
		huh.NewGroup(
			huh.NewInput().Title("Eleventh"),
			huh.NewInput().Title("Twelveth"),
			huh.NewInput().Title("Thirteenth"),
		),
	).WithLayout(huh.LayoutGrid(2, 2))
	form.Run()
}
