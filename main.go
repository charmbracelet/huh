package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
)

func main() {
	form := huh.NewForm(
		// Prompt the user for their shell hardness and base filling.
		huh.Group(
			huh.Select().
				Title("Shell?").
				Options("Hard", "Soft"),
			huh.Select().
				Title("Base").
				Options("Chicken", "Beef", "Fish", "Beans"),
		),

		// Ask which toppings they'd like and for any special instructions.
		// The customer can ask for up to 4 toppings.
		huh.Group(
			huh.MultiSelect().
				Title("Toppings").
				Options("Lettuce", "Tomatoes", "Corn", "Salsa", "Sour Cream", "Cheese").
				Filterable(true).
				Limit(4),
			huh.Text().
				Title("Special Instructions").
				CharLimit(400),
		),

		// Get some final details from the customer.
		huh.Group(
			huh.Input().
				Key("name").
				Title("What's your name?").
				Validate(huh.ValidateLength(0, 20)),
			huh.Confirm().
				Key("discount").
				Title("Would you like 15% off"),
		),
	)

	r, err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("A %s shell filled with %s and %s, topped with %s.",
		r["Shell?"], r["Base"], r["Toppings"], r["What's your name?"])

	fmt.Println("That will be $%.2f. Thanks for your order, %s!", calculatePrice(r), r["name"])
}
