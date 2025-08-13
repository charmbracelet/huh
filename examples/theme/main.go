package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh/v2"
)

func main() {
	var base *huh.Theme = huh.ThemeBase()
	var dracula *huh.Theme = huh.ThemeDracula()
	var base16 *huh.Theme = huh.ThemeBase16()
	var charm *huh.Theme = huh.ThemeCharm()
	var catppuccin *huh.Theme = huh.ThemeCatppuccin()
	var exit *huh.Theme = nil

	var theme *huh.Theme = base16

	repeat := true

	for {
		err := huh.NewSelect[*huh.Theme]().
			Title("Theme").
			Value(&theme).
			Options(
				huh.NewOption("Default", base),
				huh.NewOption("Dracula", dracula),
				huh.NewOption("Base 16", base16),
				huh.NewOption("Charm", charm),
				huh.NewOption("Catppuccin", catppuccin),
				huh.NewOption("Exit", exit),
			).Run()

		if err != nil {
			if err == huh.ErrUserAborted {
				os.Exit(130)
			}
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
				huh.NewSelect[string]().Options(huh.NewOptions("A", "B", "C")...).Title("Colors"),
				huh.NewFilePicker().Title("File"),
				huh.NewMultiSelect[string]().Options(huh.NewOptions("Red", "Green", "Yellow")...).Title("Letters"),
				huh.NewConfirm().Title("Again?").Description("Try another theme").Value(&repeat),
			),
		).WithTheme(theme).Run()
		if err != nil {
			if err == huh.ErrUserAborted {
				os.Exit(130)
			}
			fmt.Println(err)
			os.Exit(1)
		}

		if !repeat {
			break
		}
	}
}
