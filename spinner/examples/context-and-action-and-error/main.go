package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/huh/spinner"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := spinner.New().
		Context(ctx).
		ActionWithErr(func(context.Context) error {
			time.Sleep(5 * time.Second)
			return nil
		}).
		Accessible(false).
		Run()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Done!")
}
