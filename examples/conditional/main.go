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

				huh.NewSelect[string]().
					Value(&choice).
					Height(7).
					TitleFunc(func() string {
						return fmt.Sprintf("Okay, what kind of %s are you in the mood for?", category)
					}, &category).
					OptionsFunc(func() []huh.Option[string] {
						switch category {
						case fruits:
							return []huh.Option[string]{
								huh.NewOption("Tangerine", "tangerine"),
								huh.NewOption("Canteloupe", "canteloupe"),
								huh.NewOption("Pomelo", "pomelo"),
								huh.NewOption("Grapefruit", "grapefruit"),
							}
						case vegetables:
							return []huh.Option[string]{
								huh.NewOption("Carrot", "carrot"),
								huh.NewOption("Jicama", "jicama"),
								huh.NewOption("Kohlrabi", "kohlrabi"),
								huh.NewOption("Fennel", "fennel"),
								huh.NewOption("Ginger", "ginger"),
							}
						default:
							return []huh.Option[string]{
								huh.NewOption("Coffee", "coffee"),
								huh.NewOption("Tea", "tea"),
								huh.NewOption("Bubble Tea", "bubble tea"),
								huh.NewOption("Agua Fresca", "agua-fresca"),
							}
						}
					}, &category),
			),
		).Run()

	if err != nil {
		fmt.Println("Trouble in food paradise:", err)
		os.Exit(1)
	}

	fmt.Printf("One %s coming right up!\n", choice)
}
