package main

import "github.com/charmbracelet/huh"

func main() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("First"),
			huh.NewInput().Title("Second"),
			huh.NewInput().Title("Third").Description("maoe"),
		).Title("Group Title").Description("Group Description"),
	)
	form.Run()
}
