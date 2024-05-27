package main

import (
	"github.com/charmbracelet/huh"
)

func main() {
	var types []string

	form := huh.NewMultiSelect[string]().
		Filterable(false).
		Options(
			huh.NewOption("Module", "Module").Selected(true),
			huh.NewOption("Constant", "Constant").Selected(true),
			huh.NewOption("Method", "Method").Selected(true),
			huh.NewOption("Attribute", "Attribute"),
			huh.NewOption("Class", "Class").Selected(true),
		).
		Title("Types").
		Value(&types)

	huh.NewForm(huh.NewGroup(form)).Run()

}
