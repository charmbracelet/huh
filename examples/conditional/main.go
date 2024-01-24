package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
)

type consumable int

const (
	fruits consumable = iota
	vegetables
	drinks
)

func (c consumable) String() string {
	return [...]string{"fruit", "vegetable", "drink"}[c]
}

func main() {

	var category consumable
	type opts []huh.Option[string]

	var choice string

	// Then ask for a specific food item based on the previous answer.
	err :=
		huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[consumable]().
					Title("What are you in the mood for?").
					Value(&category).
					Options(
						huh.NewOption("Some fruit", fruits),
						huh.NewOption("A vegetable", vegetables),
						huh.NewOption("A drink", drinks),
					),
			),

			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Okay, what kind of fruit are you in the mood for?").
					Options(
						huh.NewOption("Tangerine", "tangerine"),
						huh.NewOption("Canteloupe", "canteloupe"),
						huh.NewOption("Pomelo", "pomelo"),
						huh.NewOption("Grapefruit", "grapefruit"),
					).
					Value(&choice),
			).WithHideFunc(func() bool { return category != fruits }),

			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Okay, what kind of vegetable are you in the mood for?").
					Options(
						huh.NewOption("Carrot", "carrot"),
						huh.NewOption("Jicama", "jicama"),
						huh.NewOption("Kohlrabi", "kohlrabi"),
						huh.NewOption("Fennel", "fennel"),
						huh.NewOption("Ginger", "ginger"),
					).
					Value(&choice),
			).WithHideFunc(func() bool { return category != vegetables }),

			huh.NewGroup(
				huh.NewSelect[string]().
					Title(fmt.Sprintf("Okay, what kind of %s are you in the mood for?", category)).
					Options(
						huh.NewOption("Coffee", "coffee"),
						huh.NewOption("Tea", "tea"),
						huh.NewOption("Bubble Tea", "bubble tea"),
						huh.NewOption("Agua Fresca", "agua-fresca"),
					).
					Value(&choice),
			).WithHideFunc(func() bool { return category != drinks }),
		).Run()

	if err != nil {
		fmt.Println("Trouble in food paradise:", err)
		os.Exit(1)
	}

	fmt.Printf("One %s coming right up!\n", choice)
}
