package main

import (
	"fmt"
	"log"

	"charm.land/bubbles/v2/key"
	"charm.land/huh/v2"
)

func main() {
	// Start with the default keymap and customize specific bindings.
	keymap := huh.NewDefaultKeyMap()

	// Change the quit key from ctrl+c to escape.
	keymap.Quit = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "quit"),
	)

	// Change select navigation to use only arrow keys (no j/k).
	keymap.Select.Up = key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "up"),
	)
	keymap.Select.Down = key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "down"),
	)

	// Disable filtering on select fields.
	keymap.Select.Filter.SetEnabled(false)

	// Change the confirm toggle to use tab instead of arrow keys.
	keymap.Confirm.Toggle = key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "toggle"),
	)

	var name string
	var color string
	var confirm bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Value(&name),
			huh.NewSelect[string]().
				Title("Favorite color").
				Options(
					huh.NewOption("Red", "red"),
					huh.NewOption("Green", "green"),
					huh.NewOption("Blue", "blue"),
				).
				Value(&color),
			huh.NewConfirm().
				Title("Confirm?").
				Value(&confirm),
		),
	).WithKeyMap(keymap)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Name: %s, Color: %s, Confirmed: %v\n", name, color, confirm)
}
