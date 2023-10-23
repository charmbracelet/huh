package huh

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"
)

// SelectStyle is the style of the select field.
type SelectStyle struct {
	Base        lipgloss.Style
	Title       lipgloss.Style
	Error       lipgloss.Style
	Description lipgloss.Style
	Cursor      lipgloss.Style
	Selected    lipgloss.Style
	Unselected  lipgloss.Style
}

// DefaultSelectStyles returns the default focused style of the select field.
func DefaultSelectStyles() (SelectStyle, SelectStyle) {
	focused := SelectStyle{
		Base:        lipgloss.NewStyle().Border(lipgloss.ThickBorder(), false).BorderLeft(true).PaddingLeft(1).MarginBottom(1).BorderForeground(lipgloss.Color("8")),
		Title:       lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		Error:       lipgloss.NewStyle().Foreground(lipgloss.Color("9")),
		Description: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Cursor:      lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		Selected:    lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
		Unselected:  lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
	}
	blurred := SelectStyle{
		Base:        lipgloss.NewStyle().Border(lipgloss.HiddenBorder(), false).BorderLeft(true).PaddingLeft(1).MarginBottom(1),
		Title:       lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Error:       lipgloss.NewStyle().Foreground(lipgloss.Color("9")),
		Description: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Cursor:      lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Selected:    lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Unselected:  lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
	}
	return focused, blurred
}

// Select is a form select field.
type Select[T any] struct {
	value        *T
	title        string
	description  string
	required     bool
	options      []Option[T]
	selected     int
	cursor       string
	err          error
	validate     func(T) error
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
		value:        new(T),
		options:      opts,
		cursor:       "> ",
		validate:     func(T) error { return nil },
		style:        &blurred,
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

// Description sets the description of the select field.
func (s *Select[T]) Description(description string) *Select[T] {
	s.description = description
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
	s.style = &s.focusedStyle
	return nil
}

// Blur blurs the select field.
func (s *Select[T]) Blur() tea.Cmd {
	s.style = &s.blurredStyle
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
	var sb strings.Builder
	sb.WriteString(s.style.Title.Render(s.title))
	if s.err != nil {
		sb.WriteString(s.style.Error.Render(" * "))
	}
	sb.WriteString("\n")
	if s.description != "" {
		sb.WriteString(s.style.Description.Render(s.description) + "\n")
	}
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
	var sb strings.Builder

	sb.WriteString(s.style.Title.Render(s.title))
	sb.WriteString("\n")

	for i, option := range s.options {
		sb.WriteString(fmt.Sprintf("%d. %s", i+1, option.Key))
		if i < len(s.options)-1 {
			sb.WriteString("\n")
		}
	}

	fmt.Println(s.style.Base.Render(sb.String()))

	option := s.options[accessibility.PromptInt("Choose: ", 1, len(s.options))-1]
	fmt.Printf("Chose: %s\n\n", option.Key)
	*s.value = option.Value
}
