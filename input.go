package huh

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Input is a form input field.
type Input struct {
	value     *string
	title     string
	required  bool
	charlimit int
	textinput textinput.Model
}

// NewInput returns a new input field.
func NewInput() *Input {
	return &Input{}
}

// Value sets the value of the input field.
func (s *Input) Value(value *string) *Input {
	s.value = value
	return s
}

// Title sets the title of the input field.
func (s *Input) Title(title string) *Input {
	s.title = title
	return s
}

// Required sets the input field as required.
func (s *Input) Required(required bool) *Input {
	s.required = required
	return s
}

// CharLimit sets the character limit of the input field.
func (s *Input) CharLimit(charlimit int) *Input {
	s.charlimit = charlimit
	return s
}

// Init initializes the input field.
func (i *Input) Init() tea.Cmd {
	i.textinput = textinput.New()
	return i.textinput.Focus()
}

// Update updates the input field.
func (i *Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	i.textinput, cmd = i.textinput.Update(msg)
	cmds = append(cmds, cmd)
	*i.value = i.textinput.Value()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			cmds = append(cmds, nextField)
		}
	}

	return i, tea.Batch(cmds...)
}

// View renders the input field.
func (i *Input) View() string {
	return i.textinput.View()
}
