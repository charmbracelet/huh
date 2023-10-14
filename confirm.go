package huh

import tea "github.com/charmbracelet/bubbletea"

// Confirm is a form confirm field.
type Confirm struct {
	value    *bool
	title    string
	required bool
}

// NewConfirm returns a new confirm field.
func NewConfirm() *Confirm {
	return &Confirm{}
}

// Value sets the value of the confirm field.
func (c *Confirm) Value(value *bool) *Confirm {
	c.value = value
	return c
}

// Title sets the title of the confirm field.
func (c *Confirm) Title(title string) *Confirm {
	c.title = title
	return c
}

// Required sets the confirm field as required.
func (c *Confirm) Required(required bool) *Confirm {
	c.required = required
	return c
}

// Init initializes the confirm field.
func (c *Confirm) Init() tea.Cmd {
	return nil
}

// Update updates the confirm field.
func (c *Confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			*c.value = true
		case "n", "N":
			*c.value = false
		case "h", "l", "left", "right":
			*c.value = !*c.value
		case "enter":
			cmds = append(cmds, nextField)
		}
	}

	return c, tea.Batch(cmds...)
}

// View renders the confirm field.
func (c *Confirm) View() string {
	return " YES " + " NO "
}
