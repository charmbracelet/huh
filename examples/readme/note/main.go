package main

import "charm.land/huh/v2"

func main() {
	note := huh.NewNote().Description(
		"# Heading\n" + "This is _italic_, *bold*" +
			"\n\n# Heading\n" + "`This is _italic_, *bold*`",
	)
	huh.NewForm(
		huh.NewGroup(note),
	).Run()
}
