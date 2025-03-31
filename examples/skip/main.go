package main

import (
	"github.com/charmbracelet/huh/v2"
)

func main() {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Charmburger").
				Description("Welcome to _Charmburger™_."),

			huh.NewSelect[string]().
				Options(huh.NewOptions("Charmburger Classic", "Chickwich", "Fishburger", "Charmpossible™ Burger")...).
				Title("Choose your burger").
				Description("At Charm we truly have a burger for everyone."),

			huh.NewNote().
				Title("🍔"),
		),

		huh.NewGroup(
			huh.NewNote().
				Title("Buy 1 get 1 free").
				Description("Welcome back to _Charmburger™_."),

			huh.NewSelect[string]().
				Options(huh.NewOptions("Charmburger Classic", "Chickwich", "Fishburger", "Charmpossible™ Burger")...).
				Title("Choose your burger").
				Description("At Charm we truly have a burger for everyone."),

			huh.NewNote().
				Title("🍔"),
		),
	)

	f.Run()
}
