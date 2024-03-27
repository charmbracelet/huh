package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/huh/spinner"
)

func main() {
	action := func(p spinner.Printer) {
		p.Println("Added bottom bun")
		time.Sleep(time.Second)
		p.Println("Added patty")
		time.Sleep(time.Second)
		p.Println("Added condiments")
		time.Sleep(time.Second)
		p.Println("Added top bun")
		time.Sleep(time.Second)
	}
	_ = spinner.New().Title("Preparing your burger").ActionWithPrinter(action).Run()
	fmt.Println("Order up!")
}
