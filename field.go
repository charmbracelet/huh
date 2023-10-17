package huh

import tea "github.com/charmbracelet/bubbletea"

// Field is a form field.
type Field interface {
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string
	Focus() tea.Cmd
	Blur() tea.Cmd
	RunAccessible()
}
