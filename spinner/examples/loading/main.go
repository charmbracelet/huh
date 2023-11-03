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
	_ = spinner.New().Title("Making your taco...").Action(action).Run()
	fmt.Println("Order up!")
}
