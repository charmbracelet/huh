package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
)

func main() {
	var name string

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What's your name?").
				Placeholder("Frank").
				Value(&name),
			huh.NewNote().
				TitleFunc(func() string {
					if name == "" {
						return "Hello!"
					}
					return fmt.Sprintf("Hello, %s!", name)
				}, &name).
				DescriptionFunc(func() string {
					if name == "" {
						return "How are you?"
					}
					return fmt.Sprintf("Your name is %d characters long", len(name))
				}, &name),
			huh.NewText().
				Title("Biography.").
				PlaceholderFunc(func() string {
					placeholder := "Tell me about yourself"
					if name != "" {
						placeholder += ", " + name
					}
					placeholder += "."
					return placeholder
				}, &name),
			huh.NewConfirm().
				TitleFunc(func() string {
					if name == "" {
						return "Continue?"
					}
					return fmt.Sprintf("Continue, %s?", name)
				}, &name).
				DescriptionFunc(func() string {
					if name == "" {
						return "Are you sure?"
					}
					return fmt.Sprintf("Last chance, %s.", name)
				}, &name),
		),
	).Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Until next time, " + name + "!")
}
