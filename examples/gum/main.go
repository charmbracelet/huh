package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("gum <input | text>")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "input":
		huh.NewInput().Run()
	case "text":
		huh.NewText().Run()
	case "confirm":
		huh.NewConfirm().Run()
	case "select":
		huh.NewSelect(os.Args[2:]...).Run()
	case "multiselect":
		huh.NewMultiSelect(os.Args[2:]...).Run()
	}
}
