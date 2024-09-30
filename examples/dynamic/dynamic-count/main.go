package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/charmbracelet/huh/v2"
)

func main() {
	var value string
	defaultValue := 10
	var chosen int

	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(&value).
				Title("Max").
				Placeholder(strconv.Itoa(defaultValue)).
				Validate(func(s string) error {
					v, err := strconv.Atoi(value)
					if err != nil {
						return err
					}
					if v <= 0 {
						return errors.New("maximum must be positive")
					}
					return nil
				}).
				Description("Select a maximum"),

			huh.NewSelect[int]().
				Value(&chosen).
				Title("Pick a number").
				DescriptionFunc(func() string {
					v, err := strconv.Atoi(value)
					if err != nil || v <= 0 {
						v = defaultValue
					}
					return "Between 1 and " + strconv.Itoa(v)
				}, &value).
				OptionsFunc(func() []huh.Option[int] {
					var options []huh.Option[int]
					v, err := strconv.Atoi(value)
					if err != nil {
						v = defaultValue
					}
					for i := range v {
						options = append(options, huh.NewOption(strconv.Itoa(i+1), i+1))
					}
					return options
				}, &value),
		),
	)
	err := f.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(chosen)
}
