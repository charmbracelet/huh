package main

import (
	"context"
	"log"
	"time"

	"charm.land/huh/v2/spinner"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/2)
	defer cancel()

	err := spinner.New().
		Context(ctx).
		WithAccessible(true).
		Run()
	if err != nil {
		log.Fatalln(err)
	}
}
