package main

import (
	"log"

	"github.com/charmbracelet/huh"
)

func main() {
	var base *huh.Theme = huh.NewBaseTheme()
	var dracula *huh.Theme = huh.NewDraculaTheme()
	var base16 *huh.Theme = huh.NewBase16Theme()
	var charm *huh.Theme = huh.NewCharmTheme()

	var theme *huh.Theme = base16

	repeat := true

	for {
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[*huh.Theme]().
					Title("Theme").
					Value(&theme).
					Options(
						huh.NewOption("Default", base),
						huh.NewOption("Dracula", dracula),
						huh.NewOption("Base 16", base16),
						huh.NewOption("Charm", charm),
					),
			),
		).Run()

		// Display form with selected theme.
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Thoughts"),
				huh.NewSelect("A", "B", "C").Title("Colors"),
				huh.NewMultiSelect("Red", "Green", "Yellow").Title("Letters"),
				huh.NewConfirm().Title("Again?").Value(&repeat),
			),
		).Theme(theme).Run()
		if err != nil {
			log.Fatal(err)
		}

		if !repeat {
			break
		}
	}
}
