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
							for _, f := range formState.Ana {
								if f == Apple {
									return true
								}
							}

							return false
						}),
						huh.NewOption(string(Banana), Banana).WithHideFunc(func() bool {
							for _, f := range formState.Ana {
								if f == Banana {
									return true
								}
							}

							return false
						}),
						huh.NewOption(string(Canteloupe), Canteloupe).WithHideFunc(func() bool {
							for _, f := range formState.Ana {
								if f == Canteloupe {
									return true
								}
							}

							return false
						}),
						huh.NewOption(string(Cherry), Cherry).WithHideFunc(func() bool {
							for _, f := range formState.Ana {
								if f == Cherry {
									return true
								}
							}

							return false
						}),
						huh.NewOption(string(Grapefruit), Grapefruit).WithHideFunc(func() bool {
							for _, f := range formState.Ana {
								if f == Grapefruit {
									return true
								}
							}

							return false
						}),
						huh.NewOption(string(Orange), Orange).WithHideFunc(func() bool {
							for _, f := range formState.Ana {
								if f == Orange {
									return true
								}
							}

							return false
						}),
						huh.NewOption(string(Pomelo), Pomelo).WithHideFunc(func() bool {
							for _, f := range formState.Ana {
								if f == Pomelo {
									return true
								}
							}

							return false
						}),
						huh.NewOption(string(Tangerine), Tangerine).WithHideFunc(func() bool {
							for _, f := range formState.Ana {
								if f == Tangerine {
									return true
								}
							}

							return false
						}),
					),
			),
			huh.NewGroup(
				huh.NewMultiSelect[Fruit]().
					Title("Which fruits Tom will choose?").
					Value(&formState.Tom).
					Options(
						huh.NewOption(string(Apple), Apple).WithHideFunc(func() bool {
							for _, f := range formState.Ana {
								if f == Apple {
									return true
								}
							}
							for _, f := range formState.Bob {
								if f == Apple {
									return true
								}
							}
							return false
						}),
						huh.NewOption(string(Banana), Banana).WithHideFunc(func() bool {
							for _, f := range formState.Ana {
								if f == Banana {
									return true
								}
							}
							for _, f := range formState.Bob {
								if f == Banana {
									return true
								}
							}
							return false
						}),
						huh.NewOption(string(Canteloupe), Canteloupe).WithHideFunc(func() bool {
							for _, f := range formState.Ana {
								if f == Canteloupe {
									return true
								}
							}
							for _, f := range formState.Bob {
								if f == Canteloupe {
									return true
								}
							}
							return false
						}),
						huh.NewOption(string(Cherry), Cherry).WithHideFunc(func() bool {
							for _, f := range formState.Ana {
								if f == Cherry {
									return true
								}
							}
							for _, f := range formState.Bob {
								if f == Cherry {
									return true
								}
							}
							return false
						}),
						huh.NewOption(string(Grapefruit), Grapefruit).WithHideFunc(func() bool {
							for _, f := range formState.Ana {
								if f == Grapefruit {
									return true
								}
							}
							for _, f := range formState.Bob {
								if f == Grapefruit {
									return true
								}
							}
							return false
						}),
						huh.NewOption(string(Orange), Orange).WithHideFunc(func() bool {
							for _, f := range formState.Ana {
								if f == Orange {
									return true
								}
							}
							for _, f := range formState.Bob {
								if f == Orange {
									return true
								}
							}
							return false
						}),
						huh.NewOption(string(Pomelo), Pomelo).WithHideFunc(func() bool {
							for _, f := range formState.Ana {
								if f == Pomelo {
									return true
								}
							}
							for _, f := range formState.Bob {
								if f == Pomelo {
									return true
								}
							}
							return false
						}),
						huh.NewOption(string(Tangerine), Tangerine).WithHideFunc(func() bool {
							for _, f := range formState.Ana {
								if f == Tangerine {
									return true
								}
							}
							for _, f := range formState.Bob {
								if f == Tangerine {
									return true
								}
							}
							return false
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
