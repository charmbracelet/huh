package main

import (
	"fmt"

	"charm.land/huh/v2"
)

func main() {
	var burger string

	err := huh.NewForm(
		// A lone skippable note in its own group is auto-skipped on start.
		huh.NewGroup(
			huh.NewNote().
				Title("Welcome").
				Description("This note is skipped automatically so the form opens on the next group."),
		),

		huh.NewGroup(
			huh.NewSelect[string]().
				Options(huh.NewOptions("Charmburger Classic", "Chickwich", "Fishburger")...).
				Title("Choose your burger").
				Value(&burger),
		),
	).Run()
	if err != nil {
		panic(err)
	}

	fmt.Printf("You chose %q\n", burger)
}
