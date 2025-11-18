package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
)

func main() {
	var name string
	form := huh.NewForm(
		huh.NewGroup(huh.NewInput().Description("What should we call you?").Value(&name)),
	).WithViewHook(func(v tea.View) tea.View {
		v.AltScreen = true
		return v
	})

	err := form.Run()
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println("Welcome, " + name + "!")
}
