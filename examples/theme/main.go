package main

import (
	"fmt"
	"os"

	"charm.land/huh/v2"
)

var themes = map[string]huh.Theme{
	"default":    huh.ThemeFunc(huh.ThemeBase),
	"dracula":    huh.ThemeFunc(huh.ThemeDracula),
	"base16":     huh.ThemeFunc(huh.ThemeBase16),
	"charm":      huh.ThemeFunc(huh.ThemeCharm),
	"catppuccin": huh.ThemeFunc(huh.ThemeCatppuccin),
}

func main() {
	theme := "base16"
	repeat := true

	for {
		err := huh.NewSelect[string]().
			Title("Theme").
			Value(&theme).
			Options(
				huh.NewOption("Default", "default"),
				huh.NewOption("Dracula", "dracula"),
				huh.NewOption("Base 16", "base16"),
				huh.NewOption("Charm", "charm"),
				huh.NewOption("Catppuccin", "catppuccin"),
				huh.NewOption("Exit", ""),
			).Run()
		if err != nil {
			if err == huh.ErrUserAborted {
				os.Exit(130)
			}
			fmt.Println(err)
			os.Exit(1)
		}
		if theme == "" {
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
		).WithTheme(themes[theme]).Run()
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
