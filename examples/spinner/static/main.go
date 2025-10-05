package main

import (
	"fmt"

	"github.com/charmbracelet/huh/v2/spinner"
)

func main() {
	_ = spinner.New().Title("Loading").WithAccessible(true).Run()
	fmt.Println("Done!")
}
