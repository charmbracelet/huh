package huh

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"
)

// MultiSelect is a form multi-select field.
type MultiSelect[T any] struct {
	title       string
	description string

	filterable bool
	limit      int

	validate func([]T) error
	err      error

	cursor int

	selected []bool
	options  []Option[T]
	value    *[]T

	focused    bool
	accessible bool
	theme      *Theme
	keymap     *MultiSelectKeyMap
}

// NewMultiSelect returns a new multi-select field.
func NewMultiSelect[T any](options ...T) *MultiSelect[T] {
	var opts []Option[T]
	for _, o := range options {
		opts = append(opts, Option[T]{Key: fmt.Sprint(o), Value: o})
	}

	return &MultiSelect[T]{
		value:    new([]T),
		options:  opts,
		selected: make([]bool, len(opts)),
		validate: func([]T) error { return nil },
		limit:    len(opts),
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

// Limit sets the limit of the multi-select field.
func (m *MultiSelect[T]) Limit(limit int) *MultiSelect[T] {
	m.limit = limit
	return m
}

// Validate sets the validation function of the multi-select field.
func (m *MultiSelect[T]) Validate(validate func([]T) error) *MultiSelect[T] {
	m.validate = validate
	return m
}

// Error returns the error of the multi-select field.
func (m *MultiSelect[T]) Error() error {
	return m.err
}

// Focus focuses the multi-select field.
func (m *MultiSelect[T]) Focus() tea.Cmd {
	m.focused = true
	return nil
}

// Blur blurs the multi-select field.
func (m *MultiSelect[T]) Blur() tea.Cmd {
	m.focused = false
	return nil
}

// KeyMap sets the keymap of the multi-select field.
func (m *MultiSelect[T]) KeyMap(k *KeyMap) Field {
	m.keymap = &k.MultiSelect
	return m
}

// KeyBinds returns the help message for the multi-select field.
func (m *MultiSelect[T]) KeyBinds() []key.Binding {
	return []key.Binding{m.keymap.Toggle, m.keymap.Up, m.keymap.Down, m.keymap.Next, m.keymap.Prev}
}

// Accessible sets the accessible mode of the multi-select field.
func (m *MultiSelect[T]) Accessible(accessible bool) Field {
	m.accessible = accessible
	return m
}

// Init initializes the multi-select field.
func (m *MultiSelect[T]) Init() tea.Cmd {
	return nil
}

// Update updates the multi-select field.
func (m *MultiSelect[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		m.err = nil

		switch {
		case key.Matches(msg, m.keymap.Up):
			m.cursor = max(m.cursor-1, 0)
		case key.Matches(msg, m.keymap.Down):
			m.cursor = min(m.cursor+1, len(m.options)-1)
		case key.Matches(msg, m.keymap.Toggle):
			if !m.selected[m.cursor] && m.numSelected() >= m.limit {
				break
			}
			m.selected[m.cursor] = !m.selected[m.cursor]
		case key.Matches(msg, m.keymap.Prev):
			m.finalize()
			if m.err != nil {
				return m, nil
			}
			return m, prevField
		case key.Matches(msg, m.keymap.Next):
			m.finalize()
			if m.err != nil {
				return m, nil
			}
			return m, nextField
		}
	}

	return m, nil
}

func (m *MultiSelect[T]) numSelected() int {
	var count int
	for _, v := range m.selected {
		if v {
			count++
		}
	}
	return count
}

func (m *MultiSelect[T]) finalize() {
	*m.value = make([]T, 0)
	for i, option := range m.options {
		if m.selected[i] {
			*m.value = append(*m.value, option.Value)
		}
	}
	m.err = m.validate(*m.value)
}

// View renders the multi-select field.
func (m *MultiSelect[T]) View() string {
	styles := m.theme.Blurred
	if m.focused {
		styles = m.theme.Focused
	}

	var sb strings.Builder
	sb.WriteString(styles.Title.Render(m.title))
	if m.err != nil {
		sb.WriteString(styles.ErrorIndicator.String())
	}
	sb.WriteString("\n")
	if m.description != "" {
		sb.WriteString(styles.Description.Render(m.description) + "\n")
	}
	c := styles.MultiSelectSelector.String()
	for i, option := range m.options {
		if m.cursor == i {
			sb.WriteString(c)
		} else {
			sb.WriteString(strings.Repeat(" ", lipgloss.Width(c)))
		}

		if m.selected[i] {
			sb.WriteString(styles.SelectedPrefix.String())
			sb.WriteString(styles.SelectedOption.Render(option.Key))
		} else {
			sb.WriteString(styles.UnselectedPrefix.String())
			sb.WriteString(styles.UnselectedOption.Render(option.Key))
		}
		if i < len(m.options)-1 {
			sb.WriteString("\n")
		}
	}
	return styles.Base.Render(sb.String())
}

func (m *MultiSelect[T]) printOptions() {
	var (
		styles = m.theme.Focused
		sb     strings.Builder
	)

	sb.WriteString(m.theme.Focused.Title.Render(m.title))
	sb.WriteString("\n")

	for i, option := range m.options {
		var prefix string
		if m.selected[i] {
			prefix = styles.SelectedPrefix.String()
			sb.WriteString(fmt.Sprintf("%d. %s%s", i+1, prefix, option.Key))
		} else {
			prefix = styles.UnselectedPrefix.String()
			sb.WriteString(fmt.Sprintf("%d. %s%s", i+1, prefix, option.Key))
		}
		if i < len(m.options)-1 {
			sb.WriteString("\n")
		}
	}

	fmt.Println(styles.Base.Render(sb.String()))
}

// Run runs the multi-select field.
func (m *MultiSelect[T]) Run() error {
	if m.accessible {
		return m.runAccessible()
	}
	return Run(m)
}

// runAccessible() runs the multi-select field in accessible mode.
func (m *MultiSelect[T]) runAccessible() error {
	styles := m.theme.Focused

	m.printOptions()

	var choice int
	for {
		fmt.Println(styles.Help.Render(fmt.Sprintf("Select up to %d options.", m.limit)))
		fmt.Println(styles.Help.Render("Type 0 to continue."))
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
	return nil
}

func (m *MultiSelect[T]) Theme(theme *Theme) Field {
	m.theme = theme
	return m
}
