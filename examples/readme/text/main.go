package main

import "github.com/charmbracelet/huh/v2"

// TODO: ensure input is not plagiarized.
func checkForPlagiarism(s string) error { return nil }

func main() {
	var story string

	text := huh.NewText().
		Title("Tell me a story.").
		Validate(checkForPlagiarism).
		Placeholder("What's on your mind?").
		Value(&story)

	// Create a form to show help.
	form := huh.NewForm(huh.NewGroup(text))
	form.Run()
}
