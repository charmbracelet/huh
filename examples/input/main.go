package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
)

func main() {
	var name string

	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Prompt("Name ").
				Value(&name),
		),
	)

	err := f.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(name)
}
