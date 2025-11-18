package main

import (
	"log"

	"charm.land/huh/v2"
	"github.com/charmbracelet/glamour/v2"
)

func main() {
	var md string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewText().Title("Markdown").Value(&md),
			huh.NewNote().Height(20).Title("Preview").
				DescriptionFunc(func() string {
					fmd, err := glamour.Render(md, "dark")
					if err != nil {
						return md
					}
					return fmd
				}, &md),
		),
	).Run()
	if err != nil {
		log.Fatal(err)
	}
}
