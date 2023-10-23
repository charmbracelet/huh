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

func main() {
	var taco Taco
	var order = Order{Taco: taco}

	// Should we run in accessible mode?
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	form := huh.NewForm(
		huh.NewGroup(huh.NewNote().
			Title("Taco Charm").
			Description("Welcome to _Taco Charm™_.\n\nHow may we take your order?").
			Next(true)),

		// What's a taco without a shell?
		// We'll need to know what filling to put inside too.
		huh.NewGroup(
			huh.NewSelect("Soft", "Hard").
				Title("Shell?").
				Description("Our shells are made fresh in-house, every day.").
				Validate(func(t string) error {
					if t == "Hard" {
						return fmt.Errorf("we're out of hard shells, sorry")
					}
					return nil
				}).
				Value(&order.Taco.Shell),

			huh.NewSelect("Chicken", "Beef", "Fish", "Beans").
				Value(&order.Taco.Base).
				Title("Base"),
		),

		// Prompt for toppings and special instructions.
		// The customer can ask for up to 4 toppings.
		huh.NewGroup(
			huh.NewMultiSelect("Lettuce", "Tomatoes", "Corn", "Salsa", "Sour Cream", "Cheese").
				Title("Toppings").
				Description("Choose up to 4.").
				Validate(func(t []string) error {
					if len(t) <= 0 {
						return fmt.Errorf("at least one topping is required")
					}
					return nil
				}).
				Value(&order.Taco.Toppings).
				Filterable(true).
				Limit(4),

			huh.NewSelect[Spice]().
				Title("Spice Level").
				Options(
					huh.NewOption("Mild", Mild),
					huh.NewOption("Medium", Medium),
					huh.NewOption("Hot", Hot),
				).
				Value(&order.Taco.Spice),
		),

		// Gather final details for the order.
		huh.NewGroup(
			huh.NewInput().
				Value(&order.Name).
				Title("What's your name?").
				Description("For when your order is ready."),

			huh.NewText().
				Value(&order.Instructions).
				Title("Special Instructions").
				CharLimit(400),

			huh.NewConfirm().
				Title("Would you like 15% off?").
				Value(&order.Discount).
				Affirmative("Yes!").
				Negative("No."),
		),
	).Accessible(accessible)

	err := form.Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("A %s shell filled with %s and topped with %s.\n", order.Taco.Shell, order.Taco.Base, strings.Join(order.Taco.Toppings, ", "))
	fmt.Printf("Thanks for your order, %s!\n", order.Name)

	if order.Discount {
		fmt.Println("Enjoy 15% off.")
	}
}
