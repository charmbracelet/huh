package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/huh/spinner"
)

func main() {
	action := func() {
		time.Sleep(2 * time.Second)
	}
	_ = spinner.New().Title("Loading").Action(action).Run()
	fmt.Println("Done!")
}
