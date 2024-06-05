package main

import (
	"github.com/charmbracelet/huh"
)

func main() {
	var file string

	huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Description("What's your name?"),

			huh.NewInput().
				Title("Username").
				Description("Select your username."),

			huh.NewFilePicker().
				Title("Profile").
				Description("Select your profile picture.").
				AllowedTypes([]string{".png", ".jpeg", ".webp", ".gif"}).
				Value(&file),

			huh.NewInput().
				Title("Password").
				EchoMode(huh.EchoModePassword).
				Description("Set your Password."),
		),
	).WithShowHelp(true).Run()
}
