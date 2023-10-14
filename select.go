package huh

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/ordered"
)

// Select is a form select field.
type Select struct {
	value    *string
	title    string
	required bool
	options  []string
	selected int
	cursor   string
}

// NewSelect returns a new select field.
func NewSelect() *Select {
	return &Select{
		cursor: " > ",
	}
}

// Value sets the value of the select field.
func (s *Select) Value(value *string) *Select {
	s.value = value
	return s
}

// Title sets the title of the select field.
func (s *Select) Title(title string) *Select {
	s.title = title
	return s
}

// Required sets the select field as required.
func (s *Select) Required(required bool) *Select {
	s.required = required
	return s
}

// Options sets the options of the select field.
func (s *Select) Options(options ...string) *Select {
	s.options = options
	return s
}

// Cursor sets the cursor of the select field.
func (s *Select) Cursor(cursor string) *Select {
	s.cursor = cursor
	return s
}

// Init initializes the select field.
func (s *Select) Init() tea.Cmd {
	return nil
}

// Update updates the select field.
func (s *Select) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			s.selected = ordered.Max(s.selected-1, 0)
		case "down", "j":
			s.selected = ordered.Min(s.selected+1, len(s.options)-1)
		case "enter":
			*s.value = s.options[s.selected]
			return s, nextField
		}
	}
	return s, nil
}

// View renders the select field.
func (s *Select) View() string {
	var sb strings.Builder
	sb.WriteString(s.title + "\n")
	for i, option := range s.options {
		if s.selected == i {
			sb.WriteString(s.cursor + option)
		} else {
			sb.WriteString(strings.Repeat(" ", lipgloss.Width(s.cursor)) + option)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
