package main

import "github.com/charmbracelet/huh/v2"

func main() {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Dynamic Help"),
			huh.NewInput().Title("Dynamic Help"),
			huh.NewInput().Title("Dynamic Help"),
		),
	)
	f.Run()
}
