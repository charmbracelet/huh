package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/huh/spinner"
)

func main() {
	action := func() {
		time.Sleep(1 * time.Second)
	}
	if err := spinner.New().Title("Preparing your burger...").Action(action).Run(); err != nil {
		fmt.Println("Failed:", err)
		return
	}
	fmt.Println("Order up!")
}
