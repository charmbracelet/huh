package main

import "github.com/charmbracelet/huh"

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
