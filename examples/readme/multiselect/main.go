package main

import "github.com/charmbracelet/huh/v2"

func main() {
	var toppings []string
	s := huh.NewMultiSelect[string]().
		Options(
			huh.NewOption("Lettuce", "Lettuce").Selected(true),
			huh.NewOption("Tomatoes", "Tomatoes").Selected(true),
			huh.NewOption("Charm Sauce", "Charm Sauce"),
			huh.NewOption("Jalapeños", "Jalapeños"),
			huh.NewOption("Cheese", "Cheese"),
			huh.NewOption("Vegan Cheese", "Vegan Cheese"),
			huh.NewOption("Nutella", "Nutella"),
		).
		Title("Toppings").
		Limit(4).
		Value(&toppings)

	huh.NewForm(huh.NewGroup(s)).Run()
}
