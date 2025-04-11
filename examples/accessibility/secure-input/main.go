package main

import (
	"errors"
	"log"

	"github.com/charmbracelet/huh"
)

func validate(s string) error {
	if s == "" {
		return errors.New("Input cannot be empty")
	}
	return nil
}

func main() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Validate(validate).
				Title("Type in your Name:"),
			huh.NewInput().
				EchoMode(huh.EchoModePassword).
				Validate(validate).
				Title("Type in your password:"),
		),
	).WithAccessible(true)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
}
