package huh

import tea "github.com/charmbracelet/bubbletea"

// Run runs a single field by wrapping it within a group and a form.
func Run(field Field, opts ...tea.ProgramOption) error {
	group := NewGroup(field)
	form := NewForm(group).WithShowHelp(false)
	return form.Run(opts...)
}
