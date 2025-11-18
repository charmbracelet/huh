package main

import (
	"context"
	"fmt"
	"time"

	"charm.land/huh/v2/spinner"
)

func main() {
	action := func() { time.Sleep(5 * time.Second) }
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	go action()
	spinner.New().Context(ctx).Run()
	fmt.Println("Done!")
}
