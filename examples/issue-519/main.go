package main

import (
	"github.com/charmbracelet/huh"
)

func main() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("FieldName").
				Placeholder("Value").
				Description("This is a complicated field that requires a long winded description to fully grasp the intricacies of its ins and outs in each particular application. Ideally it should be possible to read the entirety of this description without the user's input disappearing."),
		),
	)

	form.Run()
}
