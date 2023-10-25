package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

type Action int

const (
	Cancel Action = iota
	Push
	Fork
	Skip
)

var highlight = lipgloss.NewStyle().Foreground(lipgloss.Color("#00D7D7"))

func main() {
	var action Action
	var spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))

	repo := "charmbracelet/huh"
	theme := huh.NewBase16Theme()
	theme.FieldSeparator = lipgloss.NewStyle().SetString("\n")

	f := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[Action]().
				Value(&action).
				Options(
					huh.NewOption(repo, Push),
					huh.NewOption("Create a fork of "+repo, Fork),
					huh.NewOption("Skip pushing the branch", Skip),
					huh.NewOption("Cancel", Cancel),
				).
				Title("Where should we push the 'feature' branch?"),
		),
	).WithTheme(theme)

	err := f.Run()
	if err != nil {
		log.Fatal(err)
	}

	switch action {
	case Push:
		_ = spinner.New().Title("Pushing to charmbracelet/huh").Style(spinnerStyle).Run()
		fmt.Println("Pushed to charmbracelet/huh")
	case Fork:
		fmt.Println("Creating a fork of charmbracelet/huh...")
	case Skip:
		fmt.Println("Skipping pushing the branch...")
	case Cancel:
		fmt.Println("Cancelling...")
		os.Exit(1)
	}

	fmt.Printf("Creating pull request for %s into %s in %s\n\n", highlight.Render("test"), highlight.Render("main"), repo)

	var nextAction string

	f = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Prompt("Title ").
				Inline(true),
			huh.NewText().
				Title("Body"),
		),
		huh.NewGroup(
			huh.NewSelect("Submit", "Submit as draft", "Continue in browser", "Add metadata", "Cancel").
				Title("What's next?").Value(&nextAction),
		),
	).WithTheme(theme)

	err = f.Run()
	if err != nil {
		log.Fatal(err)
	}

	if nextAction == "Submit" {
		_ = spinner.New().Title("Submitting...").Style(spinnerStyle).Run()
		fmt.Println("Pull request submitted!")
	}
}
