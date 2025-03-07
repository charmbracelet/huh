package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

func main() {
	var toppings []string
	huh.NewForm(
		huh.NewGroup(huh.NewMultiSelect[string]().
			Options(
				huh.NewOption("Lettuce", "Lettuce"),
				huh.NewOption("Tomatoes", "Tomatoes"),
				huh.NewOption("Charm Sauce", "Charm Sauce"),
				huh.NewOption("Jalapeños", "Jalapeños"),
			).
			Title("Toppings").
			Limit(4).
			Value(&toppings))).WithAccessible(true).Run()

	fmt.Println("Selected toppings:", toppings)
}
