package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

func main() {
	var name string
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Description("What is your name?").Value(&name),
		),
	).Run(); err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("Hello, " + name + "!")
}
