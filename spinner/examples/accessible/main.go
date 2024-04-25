package main

import (
	"log"
	"time"

	"github.com/charmbracelet/huh/spinner"
)

func main() {
	err := spinner.New().
		Action(func() {
			time.Sleep(time.Second)
		}).
		Run()
	if err != nil {
		log.Fatalln(err)
	}
}
