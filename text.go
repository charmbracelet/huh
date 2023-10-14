package huh

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

// Text is a form text field.
type Text struct {
	value        *string
	title        string
	required     bool
	textarea     textarea.Model
	style        *TextStyle
	focusedStyle TextStyle
	blurredStyle TextStyle
}

// NewText returns a new text field.
func NewText() *Text {
	text := textarea.New()
	text.ShowLineNumbers = false

	f, b := DefaultTextStyles()

	return &Text{
		textarea:     text,
		focusedStyle: f,
		blurredStyle: b,
	}
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

// Focus focuses the text field.
func (s *Text) Focus() tea.Cmd {
	s.style = &s.focusedStyle
	cmd := s.textarea.Focus()
	return cmd
}

// Blur blurs the text field.
func (s *Text) Blur() tea.Cmd {
	s.style = &s.blurredStyle
	s.textarea.Blur()
	return nil
}

// Init initializes the text field.
func (s *Text) Init() tea.Cmd {
	s.textarea.FocusedStyle = s.focusedStyle.Style
	s.textarea.BlurredStyle = s.blurredStyle.Style
	s.style = &s.blurredStyle
	s.textarea.Blur()
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
	var sb strings.Builder
	sb.WriteString(s.style.Title.Render(s.title))
	sb.WriteString("\n")
	sb.WriteString(s.textarea.View())

	return sb.String()
}
