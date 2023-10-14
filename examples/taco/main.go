package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/huh"
)

type Order struct {
	Taco         Taco
	Name         string
	Instructions string
	Discount     bool
}

type Taco struct {
	Shell    string
	Base     string
	Toppings []string
}

func main() {
	var taco Taco
	var order = Order{Taco: taco}

	form := huh.NewForm(
		// What's a taco without a shell?
		// We'll need to know what filling to put inside too.
		huh.NewGroup(
			huh.NewSelect().
				Value(&order.Taco.Shell).
				Title("Shell?").
				Required(true).
				Options("Hard", "Soft"),

			huh.NewSelect().
				Value(&order.Taco.Base).
				Title("Base").
				Required(true).
				Options("Chicken", "Beef", "Fish", "Beans"),
		),

		// Prompt for toppings and special instructions.
		// The customer can ask for up to 4 toppings.
		huh.NewGroup(
			huh.NewMultiSelect().
				Value(&order.Taco.Toppings).
				Title("Toppings").
				Options("Lettuce", "Tomatoes", "Corn", "Salsa", "Sour Cream", "Cheese").
				Filterable(true).
				Limit(4),

			huh.NewText().
				Value(&order.Instructions).
				Title("Special Instructions").
				CharLimit(400),
		),

		// Gather final details for the order.
		huh.NewGroup(
			huh.NewInput().
				Value(&order.Name).
				Title("What's your name?"),

			huh.NewConfirm().
				Value(&order.Discount).
				Title("Would you like 15% off"),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("A %s shell filled with %s and topped with %s.\n", order.Taco.Shell, order.Taco.Base, strings.Join(order.Taco.Toppings, ", "))
	fmt.Printf("Thanks for your order, %s!\n", order.Name)
}
