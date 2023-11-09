package main

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/huh/spinner"
)

func main() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)

	go func() {
		// Do some work.
		time.Sleep(5 * time.Second)
		cancelFunc()
	}()

	spinner.New().Context(ctx).Run()
	fmt.Println("Done!")
}
