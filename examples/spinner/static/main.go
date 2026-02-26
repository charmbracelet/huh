package main

import (
	"fmt"

	"charm.land/huh/v2/spinner"
)

func main() {
	_ = spinner.New().Title("Loading").WithAccessible(true).Run()
	fmt.Println("Done!")
}
