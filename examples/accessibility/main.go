package main

import (
	"log"

	"github.com/charmbracelet/huh/v2"
)

func main() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(huh.NewOptions("Italian", "Greek", "Indian", "Japanese", "American")...).
				Title("Favorite Cuisine?"),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Favorite Meal?").
				Placeholder("Breakfast"),
		),
	).WithAccessible(true)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
}
