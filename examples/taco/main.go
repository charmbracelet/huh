package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
)

type Spice int

const (
	Mild Spice = iota + 1
	Medium
	Hot
)

type Order struct {
	Taco         Taco
	Name         string
	Instructions string
	Discount     bool
}

type Taco struct {
	Shell    string
	Spice    Spice
	Base     string
	Toppings []string
}

var description = `
# Taco Charm

Welcome to _Taco Charm_.

How may we take your order?

`

func main() {
	var taco Taco
	var order = Order{Taco: taco}

	// Should we run in accessible mode?
	accessible, _ := strconv.ParseBool(os.Getenv("HUH_ACCESSIBLE"))

	form := huh.NewForm(
		huh.NewGroup(huh.NewNote().Body(description).Next(true)),

		// What's a taco without a shell?
		// We'll need to know what filling to put inside too.
		huh.NewGroup(
			huh.NewSelect[string]().
				Value(&order.Taco.Shell).
				Options("Soft", "Hard").
				Title("Shell?").
				Required(true),

			huh.NewSelect[string]().
				Options("Chicken", "Beef", "Fish", "Beans").
				Value(&order.Taco.Base).
				Title("Base").
				Required(true),

			huh.NewSelect[Spice]().
				Title("Spice Level").
				OptionsKV(
					huh.NewOption("Mild", Mild),
					huh.NewOption("Medium", Medium),
					huh.NewOption("Hot", Hot),
				).
				Value(&order.Taco.Spice).
				Required(true),
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
				Title("Would you like 15% off?"),
		),
	).Accessible(accessible)

	err := form.Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("A %s shell filled with %s and topped with %s.\n", order.Taco.Shell, order.Taco.Base, strings.Join(order.Taco.Toppings, ", "))
	fmt.Printf("Thanks for your order, %s!\n", order.Name)
}
