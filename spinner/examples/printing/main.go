package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/huh/spinner"
)

func main() {
	action := func(_ context.Context, w io.Writer) error {
		fmt.Fprintln(w, "Added bottom bun")
		time.Sleep(time.Second)
		fmt.Fprintln(w, "Added patty")
		time.Sleep(time.Second)
		fmt.Fprintln(w, "Added condiments")
		time.Sleep(time.Second)
		fmt.Fprintln(w, "Added top bun")
		time.Sleep(time.Second)
		return nil
	}
	_ = spinner.New().
		Title("Preparing your burger").
		ActionErr(action).
		// Accessible(true).
		Run()
	fmt.Println("Order up!")
}
