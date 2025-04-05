package main

import (
	"fmt"

	"github.com/charmbracelet/huh/v2"
)

func isFood(_ string) error {
	return nil
}

func main() {
	var lunch string

	input := huh.NewInput().
		Title("What's for lunch?").
		Prompt("? ").
		Suggestions([]string{
			"Artichoke",
			"Baking Flour",
			"Bananas",
			"Barley",
			"Bean Sprouts",
			"Bitter Melon",
			"Black Cod",
			"Blood Orange",
			"Brown Sugar",
			"Cashew Apple",
			"Cashews",
			"Cat Food",
			"Coconut Milk",
			"Cucumber",
			"Curry Paste",
			"Currywurst",
			"Dill",
			"Dragonfruit",
			"Dried Shrimp",
			"Eggs",
			"Fish Cake",
			"Furikake",
			"Garlic",
		}).
		Validate(isFood).
		Value(&lunch)

	huh.NewForm(huh.NewGroup(input)).Run()

	fmt.Printf("Yummy, %s!\n", lunch)
}
