package huh

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

// Text is a form text field.
type Text struct {
	value    *string
	title    string
	required bool
	textarea textarea.Model
}

// NewText returns a new text field.
func NewText() *Text {
	return &Text{}
}

// Value sets the value of the text field.
func (s *Text) Value(value *string) *Text {
	s.value = value
	return s
}

// Title sets the title of the text field.
func (s *Text) Title(title string) *Text {
	s.title = title
	return s
}

// Required sets the text field as required.
func (s *Text) Required(required bool) *Text {
	s.required = required
	return s
}

// CharLimit sets the character limit of the text field.
func (s *Text) CharLimit(charlimit int) *Text {
	return s
}

// Init initializes the text field.
func (s *Text) Init() tea.Cmd {
	s.textarea = textarea.New()
	s.textarea.Focus()
	return nil
}

// Update updates the text field.
func (s *Text) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	s.textarea, cmd = s.textarea.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			cmds = append(cmds, nextField)
		}
	}

	return s, tea.Batch(cmds...)
}

// View renders the text field.
func (s *Text) View() string {
	return s.textarea.View()
}
