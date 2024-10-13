package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh/v2"
)

func main() {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(
					huh.NewOption("A", "a"),
					huh.NewOption("B", "b"),
					huh.NewOption("C", "c"),
					huh.NewOption("D", "d"),
					huh.NewOption("E", "e"),
					huh.NewOption("F", "f"),
					huh.NewOption("G", "g"),
					huh.NewOption("H", "h"),
					huh.NewOption("I", "i"),
					huh.NewOption("J", "j"),
					huh.NewOption("K", "k").Selected(true),
					huh.NewOption("L", "l"),
					huh.NewOption("M", "m"),
					huh.NewOption("N", "n"),
					huh.NewOption("O", "o"),
					huh.NewOption("P", "p"),
				),
		).WithHeight(8),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Options(
					huh.NewOption("A", "a"),
					huh.NewOption("B", "b"),
					huh.NewOption("C", "c"),
					huh.NewOption("D", "d"),
					huh.NewOption("E", "e"),
					huh.NewOption("F", "f"),
					huh.NewOption("G", "g"),
					huh.NewOption("H", "h"),
					huh.NewOption("I", "i"),
					huh.NewOption("K", "k").Selected(true),
					huh.NewOption("L", "l"),
					huh.NewOption("M", "m"),
					huh.NewOption("N", "n"),
					huh.NewOption("O", "o").Selected(true),
					huh.NewOption("P", "p"),
				),
		).WithHeight(10),
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(
					huh.NewOption("A", "a"),
					huh.NewOption("B", "b"),
					huh.NewOption("C", "c"),
					huh.NewOption("D", "d"),
					huh.NewOption("E", "e"),
					huh.NewOption("F", "f"),
					huh.NewOption("G", "g"),
					huh.NewOption("H", "h"),
					huh.NewOption("I", "i"),
					huh.NewOption("J", "j"),
					huh.NewOption("K", "k").Selected(true),
					huh.NewOption("L", "l"),
					huh.NewOption("M", "m"),
					huh.NewOption("N", "n"),
					huh.NewOption("O", "o"),
					huh.NewOption("P", "p"),
				),
		).WithHeight(5),
	)

	if err := f.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v", err)
	}
}
