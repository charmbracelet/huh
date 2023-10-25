package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
)

func main() {
	var base *huh.Theme = huh.NewBaseTheme()
	var dracula *huh.Theme = huh.NewDraculaTheme()
	var base16 *huh.Theme = huh.NewBase16Theme()
	var charm *huh.Theme = huh.NewCharmTheme()
	var exit *huh.Theme = nil

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
						huh.NewOption("Exit", exit),
					),
			),
		).Run()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if theme == nil {
			break
		}

		// Display form with selected theme.
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Thoughts").Placeholder("What's on your mind?"),
				huh.NewText().Title("More Thoughts").Placeholder("What else is on your mind?"),
				huh.NewSelect("A", "B", "C").Title("Colors"),
				huh.NewMultiSelect("Red", "Green", "Yellow").Title("Letters"),
				huh.NewConfirm().Title("Again?").Description("Try another theme").Value(&repeat),
			),
		).WithTheme(theme).Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if !repeat {
			break
		}
	}
}
