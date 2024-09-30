package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/huh/v2/spinner"
)

func main() {
	action := func() {
		time.Sleep(2 * time.Second)
	}
	_ = spinner.New().Title("Preparing your burger...").Action(action).Run()
	fmt.Println("Order up!")
}
