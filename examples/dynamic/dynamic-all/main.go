package main

import (
	"log"
	"strconv"

	"github.com/charmbracelet/huh"
)

func main() {
	var value string = "Dynamic"

	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Value(&value).Title("Dynamic").Description("Dynamic"),
			huh.NewNote().
				TitleFunc(func() string { return value }, &value).
				DescriptionFunc(func() string { return value }, &value),
			huh.NewSelect[string]().
				Height(7).
				TitleFunc(func() string { return value }, &value).
				DescriptionFunc(func() string { return value }, &value).
				OptionsFunc(func() []huh.Option[string] {
					var options []huh.Option[string]
					for i := 1; i < 6; i++ {
						options = append(options, huh.NewOption(value+" "+strconv.Itoa(i), value+strconv.Itoa(i)))
					}
					return options
				}, &value),
			huh.NewMultiSelect[string]().
				Height(7).
				TitleFunc(func() string { return value }, &value).
				DescriptionFunc(func() string { return value }, &value).
				OptionsFunc(func() []huh.Option[string] {
					var options []huh.Option[string]
					for i := 1; i < 6; i++ {
						options = append(options, huh.NewOption(value+" "+strconv.Itoa(i), value+strconv.Itoa(i)))
					}
					return options
				}, &value),
			huh.NewConfirm().
				TitleFunc(func() string { return value }, &value).
				DescriptionFunc(func() string { return value }, &value),
			huh.NewText().
				TitleFunc(func() string { return value }, &value).
				DescriptionFunc(func() string { return value }, &value).
				PlaceholderFunc(func() string { return value }, &value),
		),
	)
	err := f.Run()
	if err != nil {
		log.Fatal(err)
	}
}
