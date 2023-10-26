package main

import (
	"fmt"

	"github.com/charmbracelet/huh/spinner"
)

func main() {
	_ = spinner.New().Title("Loading").Static(true).Run()
	fmt.Println("Done!")
}
