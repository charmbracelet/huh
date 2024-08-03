package main

import (
	"log"
	"time"

	"github.com/charmbracelet/huh"
)

type Model struct {
	form *huh.Form
}

var (
	work time.Duration
	rest time.Duration
)

const (
	shortWork  = 25 * time.Minute
	mediumWork = 30 * time.Minute
	longWork   = 45 * time.Minute
)

func main() {
	breaks := map[time.Duration][]huh.Option[time.Duration]{
		shortWork: {
			huh.NewOption("5 minutes", 5*time.Minute),
			huh.NewOption("10 minutes", 10*time.Minute),
			huh.NewOption("15 minutes", 15*time.Minute),
			huh.NewOption("20 minutes", 20*time.Minute),
			huh.NewOption("25 minutes", 25*time.Minute).Selected(true),
			huh.NewOption("30 minutes", 30*time.Minute),
		},
		mediumWork: {
			huh.NewOption("5 minutes", 5*time.Minute),
			huh.NewOption("10 minutes", 10*time.Minute).Selected(true),
			huh.NewOption("15 minutes", 15*time.Minute),
			huh.NewOption("20 minutes", 20*time.Minute),
			huh.NewOption("25 minutes", 25*time.Minute),
			huh.NewOption("30 minutes", 30*time.Minute),
		},
		longWork: {
			huh.NewOption("5 minutes", 5*time.Minute),
			huh.NewOption("10 minutes", 10*time.Minute),
			huh.NewOption("15 minutes", 15*time.Minute),
			huh.NewOption("20 minutes", 20*time.Minute),
			huh.NewOption("25 minutes", 25*time.Minute),
			huh.NewOption("30 minutes", 30*time.Minute).Selected(true),
		},
	}

	form := huh.NewForm(huh.NewGroup(
		huh.NewSelect[time.Duration]().
			Title("Focus Time").
			Value(&work).
			Options(
				huh.NewOption("25 minutes", shortWork),
				huh.NewOption("30 minutes", mediumWork),
				huh.NewOption("45 minutes", longWork),
			),
		huh.NewSelect[time.Duration]().
			Value(&rest).
			Title("Break Time").
			Key("break").
			Height(8).
			OptionsFunc(func() []huh.Option[time.Duration] {
				return breaks[work]
			}, &work),
	))

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
}
