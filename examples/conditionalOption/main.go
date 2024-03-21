package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
)

type Fruit string

const (
	Tangerine  Fruit = "tangerine"
	Canteloupe Fruit = "canteloupe"
	Pomelo     Fruit = "pomelo"
	Grapefruit Fruit = "grapefruit"
	Orange     Fruit = "orange"
	Apple      Fruit = "apple"
	Banana     Fruit = "banana"
	Cherry     Fruit = "cherry"
)

func FruitWasSelectedBefore(fruit Fruit, fruits []Fruit) bool {
	for _, f := range fruits {
		if f == fruit {
			return true
		}
	}
	return false
}

type FormState struct {
	Ana []Fruit
	Bob []Fruit
	Tom []Fruit
}

func NewFormState() FormState {
	return FormState{
		Ana: []Fruit{},
		Bob: []Fruit{},
		Tom: []Fruit{},
	}
}

func main() {
	formState := NewFormState()

	// Then ask for a specific food item based on the previous answer.
	err :=
		huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[Fruit]().
					Title("Which fruits Ana will choose?").
					Value(&formState.Ana).
					Options(
						huh.NewOption(string(Apple), Apple),
						huh.NewOption(string(Banana), Banana),
						huh.NewOption(string(Canteloupe), Canteloupe),
						huh.NewOption(string(Cherry), Cherry),
						huh.NewOption(string(Grapefruit), Grapefruit),
						huh.NewOption(string(Orange), Orange),
						huh.NewOption(string(Pomelo), Pomelo),
						huh.NewOption(string(Tangerine), Tangerine),
					),
			),
			huh.NewGroup(
				huh.NewMultiSelect[Fruit]().
					Title("Which fruits Bob will choose?").
					Value(&formState.Bob).
					Options(
						huh.NewOption(string(Apple), Apple).WithHideFunc(func() bool {
							return FruitWasSelectedBefore(Apple, formState.Ana)
						}),
						huh.NewOption(string(Banana), Banana).WithHideFunc(func() bool {
							return FruitWasSelectedBefore(Banana, formState.Ana)
						}),
						huh.NewOption(string(Canteloupe), Canteloupe).WithHideFunc(func() bool {
							return FruitWasSelectedBefore(Canteloupe, formState.Ana)
						}),
						huh.NewOption(string(Cherry), Cherry).WithHideFunc(func() bool {
							return FruitWasSelectedBefore(Cherry, formState.Ana)
						}),
						huh.NewOption(string(Grapefruit), Grapefruit).WithHideFunc(func() bool {
							return FruitWasSelectedBefore(Grapefruit, formState.Ana)
						}),
						huh.NewOption(string(Orange), Orange).WithHideFunc(func() bool {
							return FruitWasSelectedBefore(Orange, formState.Ana)
						}),
						huh.NewOption(string(Pomelo), Pomelo).WithHideFunc(func() bool {
							return FruitWasSelectedBefore(Pomelo, formState.Ana)
						}),
						huh.NewOption(string(Tangerine), Tangerine).WithHideFunc(func() bool {
							return FruitWasSelectedBefore(Tangerine, formState.Ana)
						}),
					),
			),
			huh.NewGroup(
				huh.NewMultiSelect[Fruit]().
					Title("Which fruits Tom will choose?").
					Value(&formState.Tom).
					Options(
						huh.NewOption(string(Apple), Apple).WithHideFunc(func() bool {
							return FruitWasSelectedBefore(Apple, formState.Ana) || FruitWasSelectedBefore(Apple, formState.Bob)
						}),
						huh.NewOption(string(Banana), Banana).WithHideFunc(func() bool {
							return FruitWasSelectedBefore(Banana, formState.Ana) || FruitWasSelectedBefore(Banana, formState.Bob)
						}),
						huh.NewOption(string(Canteloupe), Canteloupe).WithHideFunc(func() bool {
							return FruitWasSelectedBefore(Canteloupe, formState.Ana) || FruitWasSelectedBefore(Canteloupe, formState.Bob)
						}),
						huh.NewOption(string(Cherry), Cherry).WithHideFunc(func() bool {
							return FruitWasSelectedBefore(Cherry, formState.Ana) || FruitWasSelectedBefore(Cherry, formState.Bob)
						}),
						huh.NewOption(string(Grapefruit), Grapefruit).WithHideFunc(func() bool {
							return FruitWasSelectedBefore(Grapefruit, formState.Ana) || FruitWasSelectedBefore(Grapefruit, formState.Bob)
						}),
						huh.NewOption(string(Orange), Orange).WithHideFunc(func() bool {
							return FruitWasSelectedBefore(Orange, formState.Ana) || FruitWasSelectedBefore(Orange, formState.Bob)
						}),
						huh.NewOption(string(Pomelo), Pomelo).WithHideFunc(func() bool {
							return FruitWasSelectedBefore(Pomelo, formState.Ana) || FruitWasSelectedBefore(Pomelo, formState.Bob)
						}),
						huh.NewOption(string(Tangerine), Tangerine).WithHideFunc(func() bool {
							return FruitWasSelectedBefore(Tangerine, formState.Ana) || FruitWasSelectedBefore(Tangerine, formState.Bob)
						}),
					),
			),
		).Run()

	if err != nil {
		fmt.Println("Trouble choosing fruits:", err)
		os.Exit(1)
	}

	anaFruits := make([]string, len(formState.Ana))
	for i, fruit := range formState.Ana {
		anaFruits[i] = string(fruit)
	}
	bobFruits := make([]string, len(formState.Bob))
	for i, fruit := range formState.Bob {
		bobFruits[i] = string(fruit)
	}
	tomFruits := make([]string, len(formState.Tom))
	for i, fruit := range formState.Tom {
		tomFruits[i] = string(fruit)
	}

	fmt.Printf(
		"Ana will eat %s\nBob will eat %s\nTom will eat %s",
		strings.Join(anaFruits, ", "),
		strings.Join(bobFruits, ", "),
		strings.Join(tomFruits, ", "),
	)
}
