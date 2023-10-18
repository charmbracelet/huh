package huh

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"
)

// Option is a select option.
type Option[T any] struct {
	Key   string
	Value T
}

// NewOption returns a new select option.
func NewOption[T any](key string, value T) Option[T] {
	return Option[T]{Key: key, Value: value}
}

// Select is a form select field.
type Select[T any] struct {
	value        *T
	title        string
	required     bool
	options      []Option[T]
	selected     int
	cursor       string
	style        *SelectStyle
	blurredStyle SelectStyle
	focusedStyle SelectStyle
}

// NewSelect returns a new select field.
func NewSelect[T any](options ...T) *Select[T] {
	focused, blurred := DefaultSelectStyles()

	var opts []Option[T]
	for _, option := range options {
		opts = append(opts, Option[T]{Key: fmt.Sprint(option), Value: option})
	}

	return &Select[T]{
		options:      opts,
		cursor:       "> ",
		focusedStyle: focused,
		blurredStyle: blurred,
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

// Required sets the select field as required.
func (s *Select[T]) Required(required bool) *Select[T] {
	s.required = required
	return s
}

// Options sets the options of the select field.
func (s *Select[T]) Options(options ...Option[T]) *Select[T] {
	s.options = options
	return s
}

// Cursor sets the cursor of the select field.
func (s *Select[T]) Cursor(cursor string) *Select[T] {
	s.cursor = cursor
	return s
}

// Styles sets the styles of the select field.
func (s *Select[T]) Styles(focused, blurred SelectStyle) *Select[T] {
	s.blurredStyle = blurred
	s.focusedStyle = focused
	return s
}

// Focus focuses the select field.
func (s *Select[T]) Focus() tea.Cmd {
	s.style = &s.focusedStyle
	return nil
}

// Blur blurs the select field.
func (s *Select[T]) Blur() tea.Cmd {
	s.style = &s.blurredStyle
	return nil
}

// Init initializes the select field.
func (s *Select[T]) Init() tea.Cmd {
	s.style = &s.blurredStyle
	return nil
}

// Update updates the select field.
func (s *Select[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			*s.value = s.options[s.selected].Value
			return s, nextField
		}
	}
	return s, nil
}

// View renders the select field.
func (s *Select[T]) View() string {
	var sb strings.Builder
	sb.WriteString(s.style.Title.Render(s.title) + "\n")
	c := s.style.Cursor.Render(s.cursor)
	for i, option := range s.options {
		if s.selected == i {
			sb.WriteString(c + s.style.Selected.Render(option.Key))
		} else {
			sb.WriteString(strings.Repeat(" ", lipgloss.Width(c)) + s.style.Unselected.Render(option.Key))
		}
		if i < len(s.options)-1 {
			sb.WriteString("\n")
		}
	}
	return s.style.Base.Render(sb.String())
}

// Run runs an accessible select field.
func (s *Select[T]) Run() {
	fmt.Println(s.style.Title.Render(s.title))
	for i, option := range s.options {
		fmt.Printf("%d. %s\n", i+1, option)
	}

	option := s.options[accessibility.PromptInt(1, len(s.options))-1]
	fmt.Printf("Selected: %s\n\n", option)
	*s.value = option.Value
}
