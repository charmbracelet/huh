package main

import (
	"log"

	"github.com/charmbracelet/huh"
)

func main() {
	form := huh.NewForm(
		// What's a taco without a shell?
		// We'll need to know what filling to put inside too.
		huh.NewGroup(
			huh.NewSelect("Hard", "Soft").
				Title("Shell?"),

			huh.NewSelect("Chicken", "Beef", "Fish", "Beans").
				Title("Base"),
		),

		// Prompt for toppings and special instructions.
		// The customer can ask for up to 4 toppings.
		huh.NewGroup(
			huh.NewMultiSelect("Lettuce", "Tomatoes", "Corn", "Salsa", "Sour Cream", "Cheese").
				Title("Toppings").
				Limit(4),

			huh.NewText().
				Title("Special Instructions").
				CharLimit(400),
		),

		// Gather final details for the order.
		huh.NewGroup(
			huh.NewInput().
				Title("What's your name?").
				Validate(validateName),

			huh.NewConfirm().
				Title("Would you like 15% off"),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
}
