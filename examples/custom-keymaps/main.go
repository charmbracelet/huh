package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/key"
	"charm.land/huh/v2"
)

func inputKeyMap(next string) *huh.KeyMap {
	keymap := huh.NewDefaultKeyMap()
	keymap.Input.Next = key.NewBinding(key.WithKeys(next), key.WithHelp(next, "next"))
	keymap.Input.Submit = key.NewBinding(key.WithKeys(next), key.WithHelp(next, "submit"))
	return keymap
}

func main() {
	formInput := huh.NewInput().Title("Form keymap (press F1)")
	groupInput := huh.NewInput().Title("Group keymap (press F2)")
	input := huh.NewInput().Title("Input keymap (press F3)")

	group := huh.NewGroup(groupInput)
	form := huh.NewForm(
		huh.NewGroup(formInput),
		group,
		huh.NewGroup(input),
	)

	// Apply keymaps from broadest to narrowest because each call also updates its children.
	form.WithKeyMap(inputKeyMap("f1"))
	group.WithKeyMap(inputKeyMap("f2"))
	input.WithKeyMap(inputKeyMap("f3"))
	form.UpdateFieldPositions()

	if err := form.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
