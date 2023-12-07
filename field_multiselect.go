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
	value *[]T
	key   string

	// customization
	title       string
	description string
	options     []Option[T]
	filterable  bool
	limit       int

	// error handling
	validate func([]T) error
	err      error

	// state
	cursor  int
	focused bool

	// options
	width      int
	accessible bool
	theme      *Theme
	keymap     *MultiSelectKeyMap
}

// NewMultiSelect returns a new multi-select field.
func NewMultiSelect[T any]() *MultiSelect[T] {
	return &MultiSelect[T]{
		options:  []Option[T]{},
		value:    new([]T),
		validate: func([]T) error { return nil },
	}
}

// Value sets the value of the multi-select field.
func (m *MultiSelect[T]) Value(value *[]T) *MultiSelect[T] {
	m.value = value
	return m
}

// Key sets the key of the select field which can be used to retrieve the value
// after submission.
func (m *MultiSelect[T]) Key(key string) *MultiSelect[T] {
	m.key = key
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
	if len(options) <= 0 {
		return m
	}
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

// KeyBinds returns the help message for the multi-select field.
func (m *MultiSelect[T]) KeyBinds() []key.Binding {
	return []key.Binding{m.keymap.Toggle, m.keymap.Up, m.keymap.Down, m.keymap.Next, m.keymap.Prev}
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
			if !m.options[m.cursor].selected && m.limit > 0 && m.numSelected() >= m.limit {
				break
			}
			m.options[m.cursor].selected = !m.options[m.cursor].selected
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
	for _, o := range m.options {
		if o.selected {
			count++
		}
	}
	return count
}

func (m *MultiSelect[T]) finalize() {
	*m.value = make([]T, 0)
	for _, option := range m.options {
		if option.selected {
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

		if m.options[i].selected {
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
		sb strings.Builder
	)

	sb.WriteString(m.theme.Focused.Title.Render(m.title))
	sb.WriteString("\n")

	for i, option := range m.options {
		if option.selected {
			sb.WriteString(m.theme.Focused.SelectedOption.Render(fmt.Sprintf("%d. %s %s", i+1, "✓", option.Key)))
		} else {
			sb.WriteString(fmt.Sprintf("%d. %s %s", i+1, " ", option.Key))
		}
		sb.WriteString("\n")
	}

	fmt.Println(m.theme.Blurred.Base.Render(sb.String()))
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
	m.printOptions()

	var choice int
	for {
		fmt.Printf("Select up to %d options. 0 to continue.\n", m.limit)

		choice = accessibility.PromptInt("Select: ", 0, len(m.options))
		if choice == 0 {
			m.finalize()
			err := m.validate(*m.value)
			if err != nil {
				fmt.Println(err)
				continue
			}
			break
		}
		m.options[choice-1].selected = !m.options[choice-1].selected
		if m.options[choice-1].selected {
			fmt.Printf("Selected: %s\n\n", m.options[choice-1].Key)
		} else {
			fmt.Printf("Deselected: %s\n\n", m.options[choice-1].Key)
		}

		m.printOptions()
	}

	var values []string

	for _, option := range m.options {
		if option.selected {
			*m.value = append(*m.value, option.Value)
			values = append(values, option.Key)
		}
	}

	fmt.Println(m.theme.Focused.SelectedOption.Render("Selected:", strings.Join(values, ", ")+"\n"))
	return nil
}

// WithTheme sets the theme of the multi-select field.
func (m *MultiSelect[T]) WithTheme(theme *Theme) Field {
	m.theme = theme
	return m
}

// WithKeyMap sets the keymap of the multi-select field.
func (m *MultiSelect[T]) WithKeyMap(k *KeyMap) Field {
	m.keymap = &k.MultiSelect
	return m
}

// WithAccessible sets the accessible mode of the multi-select field.
func (m *MultiSelect[T]) WithAccessible(accessible bool) Field {
	m.accessible = accessible
	return m
}

// WithWidth sets the width of the multi-select field.
func (m *MultiSelect[T]) WithWidth(width int) Field {
	m.width = width
	return m
}

// GetKey returns the multi-select's key.
func (m *MultiSelect[T]) GetKey() string {
	return m.key
}

// GetValue returns the multi-select's value.
func (m *MultiSelect[T]) GetValue() any {
	return *m.value
}
