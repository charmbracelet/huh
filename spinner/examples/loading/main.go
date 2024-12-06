package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/huh/spinner"
)

func main() {
	action := func() {
		time.Sleep(2 * time.Second)
	}
	if err := spinner.New().Title("Preparing your burger...").Action(action).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Order up!")
}
