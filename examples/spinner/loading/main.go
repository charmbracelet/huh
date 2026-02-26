package main

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2/spinner"
)

func main() {
	action := func() {
		time.Sleep(1 * time.Second)
	}
	if err := spinner.New().
		Title("Preparing your burger...").
		Action(action).
		WithViewHook(func(v tea.View) tea.View {
			v.ProgressBar = tea.NewProgressBar(tea.ProgressBarIndeterminate, 1)
			return v
		}).
		Run(); err != nil {
		fmt.Println("Failed:", err)
		return
	}
	fmt.Println("Order up!")
}
