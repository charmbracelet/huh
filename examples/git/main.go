package main

import (
	"github.com/charmbracelet/huh/v2"
)

// types is the possible commit types specified by the conventional commit spec.
var types = []string{"fix", "feat", "docs", "style", "refactor", "test", "chore", "revert"}

// This form is used to write a conventional commit message. It prompts the user
// to choose the type of commit as specified in the conventional commit spec.
// And then prompts for the summary and detailed description of the message and
// uses the values provided as the summary and details of the message.
func main() {
	var commit, scope string
	var summary, description string
	var confirm bool

	huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Type").Value(&commit).Placeholder("feat").Suggestions(types),
			huh.NewInput().Title("Scope").Value(&scope).Placeholder("scope"),
		),
		huh.NewGroup(
			huh.NewInput().Title("Summary").Value(&summary).Placeholder("Summary of changes"),
			huh.NewText().Title("Description").Value(&description).Placeholder("Detailed description of changes"),
		),
		huh.NewGroup(huh.NewConfirm().Title("Commit changes?").Value(&confirm)),
	).Run()
}
