package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/huh/v2"
)

func main() {

	count := 0
	go func() {
		for {
			count++
			time.Sleep(1 * time.Second)
		}
	}()

	descriptionFunc := func() string {
		return fmt.Sprintf("The count is: %d", count)
	}

	huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Fill in the input").
			DescriptionFunc(descriptionFunc, &count),
		huh.NewInput().
			Title("Fill in the input").
			DescriptionFunc(descriptionFunc, &count),
	)).Run()

}
