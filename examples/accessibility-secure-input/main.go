package main

import (
	"errors"
	"log"

	"github.com/charmbracelet/huh"
)

func validate(s string) error {
	if s == "" {
		return errors.New("input cannot be empty")
	}
	return nil
}

func main() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Welcome!").
				Description("This is an accessible form example!"),
			huh.NewInput().
				Validate(validate).
				Title("Name:"),
			huh.NewInput().
				EchoMode(huh.EchoModePassword).
				Validate(validate).
				Title("Password:"),
			huh.NewMultiSelect[string]().
				Options(huh.NewOptions(
					"Red",
					"Green",
					"Yellow",
				)...).
				Limit(2).
				Title("Choose some colors:"),
			huh.NewSelect[string]().
				Options(huh.NewOptions(
					"Red",
					"Green",
					"Yellow",
				)...).
				Title("Choose the best color:"),
			huh.NewFilePicker().
				Title("Which file?"),
			huh.NewConfirm().
				Title("Send something?"),
		),
	).WithAccessible(true)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
}
