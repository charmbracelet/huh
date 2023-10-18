package huh

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"
)

// Select is a form select field.
type Select struct {
	value        *string
	title        string
	required     bool
	options      []string
	selected     int
	cursor       string
	style        *SelectStyle
	blurredStyle SelectStyle
	focusedStyle SelectStyle
}

// NewSelect returns a new select field.
func NewSelect() *Select {
	focused, blurred := DefaultSelectStyles()
	return &Select{
		cursor:       "> ",
		focusedStyle: focused,
		blurredStyle: blurred,
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

// Styles sets the styles of the select field.
func (s *Select) Styles(focused, blurred SelectStyle) *Select {
	s.blurredStyle = blurred
	s.focusedStyle = focused
	return s
}

// Focus focuses the select field.
func (s *Select) Focus() tea.Cmd {
	s.style = &s.focusedStyle
	return nil
}

// Blur blurs the select field.
func (s *Select) Blur() tea.Cmd {
	s.style = &s.blurredStyle
	return nil
}

// Init initializes the select field.
func (s *Select) Init() tea.Cmd {
	s.style = &s.blurredStyle
	return nil
}

// Update updates the select field.
func (s *Select) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			s.selected = max(s.selected-1, 0)
		case "down", "j":
			s.selected = min(s.selected+1, len(s.options)-1)
		case "shift+tab":
			return s, prevField
		case "tab", "enter":
			*s.value = s.options[s.selected]
			return s, nextField
		}
	}
	return s, nil
}

// View renders the select field.
func (s *Select) View() string {
	var sb strings.Builder
	sb.WriteString(s.style.Title.Render(s.title) + "\n")
	c := s.style.Cursor.Render(s.cursor)
	for i, option := range s.options {
		if s.selected == i {
			sb.WriteString(c + s.style.Selected.Render(option))
		} else {
			sb.WriteString(strings.Repeat(" ", lipgloss.Width(c)) + s.style.Unselected.Render(option))
		}
		if i < len(s.options)-1 {
			sb.WriteString("\n")
		}
	}
	return s.style.Base.Render(sb.String())
}

// Run runs an accessible select field.
func (s *Select) Run() {
	fmt.Println(s.style.Title.Render(s.title))
	for i, option := range s.options {
		fmt.Printf("%d. %s\n", i+1, option)
	}

	value := s.options[accessibility.PromptInt(1, len(s.options))-1]
	fmt.Printf("Selected: %s\n\n", value)
	*s.value = value
}
