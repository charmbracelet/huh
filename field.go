package huh

import tea "github.com/charmbracelet/bubbletea"

// Field is a form field.
type Field interface {
	// Bubble Tea Model
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string

	// Bubble Tea Events
	Blur() tea.Cmd
	Focus() tea.Cmd

	// Accessible Prompt (non-redraw)
	Run()
}
