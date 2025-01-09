package main

import "github.com/charmbracelet/huh"

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
			huh.NewInput().Title("Eighth"),
			huh.NewInput().Title("Nineth"),
			huh.NewInput().Title("Tenth"),
		),
	).WithLayout(huh.LayoutStack)
	form.Run()
}
