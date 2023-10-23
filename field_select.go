package huh

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"
)

// Select is a form select field.
type Select[T any] struct {
	value       *T
	title       string
	description string

	validate func(T) error
	err      error

	options  []Option[T]
	selected int

	focused bool
	theme   *Theme
}

// NewSelect returns a new select field.
func NewSelect[T any](options ...T) *Select[T] {
	var opts []Option[T]
	for _, option := range options {
		opts = append(opts, Option[T]{Key: fmt.Sprint(option), Value: option})
	}

	return &Select[T]{
		value:    new(T),
		options:  opts,
		validate: func(T) error { return nil },
	}
}

// Value sets the value of the select field.
func (s *Select[T]) Value(value *T) *Select[T] {
	s.value = value
	return s
}

// Title sets the title of the select field.
func (s *Select[T]) Title(title string) *Select[T] {
	s.title = title
	return s
}

// Description sets the description of the select field.
func (s *Select[T]) Description(description string) *Select[T] {
	s.description = description
	return s
}

// Options sets the options of the select field.
func (s *Select[T]) Options(options ...Option[T]) *Select[T] {
	s.options = options
	return s
}

// Validate sets the validation function of the select field.
func (s *Select[T]) Validate(validate func(T) error) *Select[T] {
	s.validate = validate
	return s
}

// Error returns the error of the select field.
func (s *Select[T]) Error() error {
	return s.err
}

// Focus focuses the select field.
func (s *Select[T]) Focus() tea.Cmd {
	s.focused = true
	return nil
}

// Blur blurs the select field.
func (s *Select[T]) Blur() tea.Cmd {
	s.focused = false
	s.err = s.validate(*s.value)
	return nil
}

// Init initializes the select field.
func (s *Select[T]) Init() tea.Cmd {
	return nil
}

// Update updates the select field.
func (s *Select[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		s.err = nil
		switch msg.String() {
		case "up", "k":
			s.selected = max(s.selected-1, 0)
		case "down", "j":
			s.selected = min(s.selected+1, len(s.options)-1)
		case "shift+tab":
			return s, prevField
		case "tab", "enter":
			value := s.options[s.selected].Value
			s.err = s.validate(value)
			if s.err != nil {
				return s, nil
			}
			*s.value = value
			return s, nextField
		}
	}
	return s, nil
}

// View renders the select field.
func (s *Select[T]) View() string {
	styles := s.theme.Blurred
	if s.focused {
		styles = s.theme.Focused
	}

	var sb strings.Builder
	sb.WriteString(styles.Title.Render(s.title))
	if s.err != nil {
		sb.WriteString(styles.Error.Render(" * "))
	}
	sb.WriteString("\n")
	if s.description != "" {
		sb.WriteString(styles.Description.Render(s.description) + "\n")
	}

	c := styles.SelectSelector.String()
	for i, option := range s.options {
		if s.selected == i {
			sb.WriteString(c + styles.Option.Render(option.Key))
		} else {
			sb.WriteString(strings.Repeat(" ", lipgloss.Width(c)) + styles.Option.Render(option.Key))
		}
		if i < len(s.options)-1 {
			sb.WriteString("\n")
		}
	}
	return styles.Base.Render(sb.String())
}

// Run runs an accessible select field.
func (s *Select[T]) Run() {
	var sb strings.Builder

	sb.WriteString(s.theme.Focused.Title.Render(s.title) + "\n")

	for i, option := range s.options {
		sb.WriteString(fmt.Sprintf("%d. %s", i+1, option.Key))
		if i < len(s.options)-1 {
			sb.WriteString("\n")
		}
	}

	fmt.Println(s.theme.Focused.Base.Render(sb.String()))

	option := s.options[accessibility.PromptInt("Choose: ", 1, len(s.options))-1]
	fmt.Printf("Chose: %s\n\n", option.Key)
	*s.value = option.Value
}

func (s *Select[T]) Theme(theme *Theme) Field {
	s.theme = theme
	return s
}
