package huh

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"
)

// MultiSelectStyle is the style of the multi-select field.
type MultiSelectStyle struct {
	Base             lipgloss.Style
	Title            lipgloss.Style
	Description      lipgloss.Style
	Help             lipgloss.Style
	Cursor           lipgloss.Style
	Selected         lipgloss.Style
	Unselected       lipgloss.Style
	SelectedPrefix   lipgloss.Style
	UnselectedPrefix lipgloss.Style
}

// DefaultMultiSelectStyles returns the default focused style of the multi-select field.
func DefaultMultiSelectStyles() (MultiSelectStyle, MultiSelectStyle) {
	focused := MultiSelectStyle{
		Base:             lipgloss.NewStyle().Border(lipgloss.ThickBorder(), false).BorderLeft(true).PaddingLeft(1).MarginBottom(1).BorderForeground(lipgloss.Color("8")),
		Title:            lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		Description:      lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Help:             lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Cursor:           lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		Selected:         lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
		Unselected:       lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
		SelectedPrefix:   lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
		UnselectedPrefix: lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
	}
	blurred := MultiSelectStyle{
		Base:             lipgloss.NewStyle().Border(lipgloss.HiddenBorder(), false).BorderLeft(true).PaddingLeft(1).MarginBottom(1),
		Title:            lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Description:      lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Help:             lipgloss.NewStyle().Foreground(lipgloss.Color("0")),
		Cursor:           lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Selected:         lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Unselected:       lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		SelectedPrefix:   lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		UnselectedPrefix: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
	}
	return focused, blurred
}

// MultiSelect is a form multi-select field.
type MultiSelect[T any] struct {
	title            string
	description      string
	required         bool
	filterable       bool
	limit            int
	cursor           int
	cursorPrefix     string
	selectedPrefix   string
	unselectedPrefix string
	selected         []bool
	options          []Option[T]
	value            *[]T
	style            *MultiSelectStyle
	blurredStyle     MultiSelectStyle
	focusedStyle     MultiSelectStyle
	theme            *Theme
}

// NewMultiSelect returns a new multi-select field.
func NewMultiSelect[T any](options ...T) *MultiSelect[T] {
	f, b := DefaultMultiSelectStyles()

	var opts []Option[T]
	for _, o := range options {
		opts = append(opts, Option[T]{Key: fmt.Sprint(o), Value: o})
	}

	return &MultiSelect[T]{
		value:            new([]T),
		options:          opts,
		cursorPrefix:     "> ",
		selectedPrefix:   "[•] ",
		unselectedPrefix: "[ ] ",
		focusedStyle:     f,
		blurredStyle:     b,
		style:            &b,
		selected:         make([]bool, len(opts)),
	}
}

// Value sets the value of the multi-select field.
func (m *MultiSelect[T]) Value(value *[]T) *MultiSelect[T] {
	m.value = value
	return m
}

// Title sets the title of the multi-select field.
func (m *MultiSelect[T]) Title(title string) *MultiSelect[T] {
	m.title = title
	return m
}

// Description sets the description of the multi-select field.
func (m *MultiSelect[T]) Description(description string) *MultiSelect[T] {
	m.description = description
	return m
}

// Required sets the multi-select field as required.
func (m *MultiSelect[T]) Required(required bool) *MultiSelect[T] {
	m.required = required
	return m
}

// Options sets the options of the multi-select field.
func (m *MultiSelect[T]) Options(options ...Option[T]) *MultiSelect[T] {
	m.options = options
	return m
}

// Filterable sets the multi-select field as filterable.
func (m *MultiSelect[T]) Filterable(filterable bool) *MultiSelect[T] {
	m.filterable = filterable
	return m
}

// Cursor sets the cursor of the multi-select field.
func (m *MultiSelect[T]) Cursor(cursor string) *MultiSelect[T] {
	m.cursorPrefix = cursor
	return m
}

// Limit sets the limit of the multi-select field.
func (m *MultiSelect[T]) Limit(limit int) *MultiSelect[T] {
	m.limit = limit
	return m
}

// Focus focuses the multi-select field.
func (m *MultiSelect[T]) Focus() tea.Cmd {
	m.style = &m.focusedStyle
	return nil
}

// Blur blurs the multi-select field.
func (m *MultiSelect[T]) Blur() tea.Cmd {
	m.style = &m.blurredStyle
	return nil
}

// Init initializes the multi-select field.
func (m *MultiSelect[T]) Init() tea.Cmd {
	return nil
}

// Update updates the multi-select field.
func (m *MultiSelect[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.cursor = max(m.cursor-1, 0)
		case "down", "j":
			m.cursor = min(m.cursor+1, len(m.options)-1)
		case " ", "x":
			m.selected[m.cursor] = !m.selected[m.cursor]
		case "shift+tab":
			m.finalize()
			return m, prevField
		case "tab", "enter":
			m.finalize()
			return m, nextField
		}
	}

	return m, nil
}

func (m *MultiSelect[T]) finalize() {
	*m.value = make([]T, 0)
	for i, option := range m.options {
		if m.selected[i] {
			*m.value = append(*m.value, option.Value)
		}
	}
}

// View renders the multi-select field.
func (m *MultiSelect[T]) View() string {
	var sb strings.Builder
	sb.WriteString(m.style.Title.Render(m.title) + "\n")
	if m.description != "" {
		sb.WriteString(m.style.Description.Render(m.description) + "\n")
	}
	c := m.style.Cursor.Render(m.cursorPrefix)
	for i, option := range m.options {
		if m.cursor == i {
			sb.WriteString(c)
		} else {
			sb.WriteString(strings.Repeat(" ", lipgloss.Width(c)))
		}

		if m.selected[i] {
			sb.WriteString(m.style.SelectedPrefix.Render(m.selectedPrefix))
			sb.WriteString(m.style.Selected.Render(option.Key))
		} else {
			sb.WriteString(m.style.UnselectedPrefix.Render(m.unselectedPrefix))
			sb.WriteString(m.style.Unselected.Render(option.Key))
		}
		if i < len(m.options)-1 {
			sb.WriteString("\n")
		}
	}
	return m.style.Base.Render(sb.String())
}

func (m *MultiSelect[T]) printOptions() {
	var sb strings.Builder

	sb.WriteString(m.style.Title.Render(m.title))
	sb.WriteString("\n")

	for i, option := range m.options {
		var prefix string
		if m.selected[i] {
			prefix = m.style.SelectedPrefix.Render(m.selectedPrefix)
			sb.WriteString(fmt.Sprintf("%d. %s%s", i+1, prefix, option.Key))
		} else {
			prefix = m.style.UnselectedPrefix.Render(m.unselectedPrefix)
			sb.WriteString(fmt.Sprintf("%d. %s%s", i+1, prefix, option.Key))
		}
		if i < len(m.options)-1 {
			sb.WriteString("\n")
		}
	}

	fmt.Println(m.style.Base.Render(sb.String()))
}

// Run runs the multi-select field in accessible mode.
func (m *MultiSelect[T]) Run() {
	m.printOptions()

	var choice int
	for {
		fmt.Println(m.style.Help.Render(fmt.Sprintf("Select up to %d options.", m.limit)))
		fmt.Println(m.style.Help.Render("Type 0 to continue."))
		fmt.Println()

		choice = accessibility.PromptInt("Select: ", 0, len(m.options))
		if choice == 0 {
			break
		}
		m.selected[choice-1] = !m.selected[choice-1]
		if m.selected[choice-1] {
			fmt.Printf("Selected: %d\n\n", choice)
		} else {
			fmt.Printf("Deselected: %d\n\n", choice)
		}

		m.printOptions()
	}

	var values []string

	for i, option := range m.options {
		if m.selected[i] {
			*m.value = append(*m.value, option.Value)
			values = append(values, option.Key)
		}
	}

	fmt.Println("Selected:", strings.Join(values, ", ")+"\n")
}

func (m *MultiSelect[T]) setTheme(theme *Theme) {
	m.theme = theme
}
