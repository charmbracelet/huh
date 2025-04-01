package main

import (
	"fmt"

	"github.com/charmbracelet/huh/v2/spinner"
)

func main() {
	_ = spinner.New().Title("Loading").Accessible(true).Run()
	fmt.Println("Done!")
}
