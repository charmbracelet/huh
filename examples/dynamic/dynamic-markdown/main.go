package main

import (
	"log"

	"github.com/charmbracelet/huh"
)

func main() {
	var md string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewText().Title("Markdown").Value(&md),
			huh.NewNote().Height(20).Title("Preview").DescriptionFunc(func() string { return md }, &md),
		),
	).Run()
	if err != nil {
		log.Fatal(err)
	}
}
