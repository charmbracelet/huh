package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	xstrings "github.com/charmbracelet/x/exp/strings"
)

func main() {
	var fruits []string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select your favorite fruits").
				Options(huh.NewOptions(
					"Apple - Macintosh",
					"Apple - Granny Smith",
					"Apple - Honeycrisp",
					"Citrus - Orange",
					"Citrus - Grapefruit",
					"Berry - Strawberry",
					"Berry - Blueberry",
					"Berry - Raspberry",
					"Berry - Blackberry",
				)...).
				Filterable(true).
				Value(&fruits),
		),
	)

	if err := form.Run(); err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	fruitCategories := make([]string, 0, 3)
	for _, f := range fruits {
		category := strings.Split(f, " - ")[0]
		add := true
		for _, c := range fruitCategories {
			if c == category {
				add = false
				break
			}
		}
		if add {
			fruitCategories = append(fruitCategories, category)
		}
	}

	var sb strings.Builder
	keyword := func(s string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
	}
	fmt.Fprintf(&sb,
		"%s\n%s",
		lipgloss.NewStyle().Bold(true).Render("Your preferred fruit categories"),
		keyword(xstrings.EnglishJoin(fruitCategories, true)),
	)

	fmt.Println(
		lipgloss.NewStyle().
			Width(40).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Render(sb.String()),
	)
}
