package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/charmbracelet/huh/spinner"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := spinner.New().
		Context(ctx).
		ActionErr(func(context.Context, io.Writer) error {
			time.Sleep(time.Minute)
			return nil
		}).
		Accessible(false).
		Run()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Done!")
}
