package main

import "github.com/charmbracelet/huh"

func main() {
	var toppings []string
	s := huh.NewMultiSelect[string]().
		Options(
			huh.NewOption("Cheese", "cheese").Selected(true),
			huh.NewOption("Lettuce", "lettuce").Selected(true),
			huh.NewOption("Corn", "corn"),
			huh.NewOption("Salsa", "salsa"),
			huh.NewOption("Sour Cream", "sour cream"),
			huh.NewOption("Tomatoes", "tomatoes"),
		).
		Title("Toppings").
		Limit(4).
		Value(&toppings)

	huh.NewForm(huh.NewGroup(s)).Run()
}
