package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/huh/spinner"
)

func main() {
	action := func() error {
		time.Sleep(1 * time.Second)
		return nil
	}
	if err := spinner.New().Title("Preparing your burger...").Action(action).Run(); err != nil {
		fmt.Println("Failed:", err)
		return
	}
	fmt.Println("Order up!")
}
