package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/huh/v2"
)

func main() {
	var name string
	form := huh.NewForm(
		huh.NewGroup(huh.NewInput().Description("What should we call you?").Value(&name)),
	).WithProgramOptions(tea.WithAltScreen())

	err := form.Run()
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println("Welcome, " + name + "!")
}
