package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh/v2"
	"github.com/charmbracelet/huh/v2/spinner"
)

func main() {
	var org string
	var repo string

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(&org).
				Title("Organization").
				Placeholder("charmbracelet"),
			huh.NewInput().
				Value(&repo).
				Title("Repository").
				PlaceholderFunc(func() string {
					switch org {
					case "hashicorp":
						return "terraform"
					case "golang":
						return "go"
					default: // charmbracelet
						return "bubbletea"
					}
				}, &org).
				SuggestionsFunc(func() []string {
					switch org {
					case "charmbracelet":
						return []string{"bubbletea", "huh", "mods", "melt", "freeze", "gum", "vhs", "pop", "lipgloss", "harmonica"}
					case "hashicorp":
						return []string{"terraform", "vault", "waypoint"}
					case "golang":
						return []string{"go", "net", "sys", "text", "tools"}
					default:
						return nil
					}
				}, &org),
		),
	).Run()
	if err != nil {
		log.Fatal(err)
	}

	spinner.New().Title(fmt.Sprintf("Cloning %s/%s...", org, repo)).Run()
}
