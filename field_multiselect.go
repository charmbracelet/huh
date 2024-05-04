package huh

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"
)

// MultiSelect is a form multi-select field.
type MultiSelect[T comparable] struct {
	value *[]T
	key   string

	// customization
	title           string
	description     string
	options         []Option[T]
	filterable      bool
	filteredOptions []Option[T]
	limit           int
	height          int

	// error handling
	validate func([]T) error
	err      error

	// state
	cursor    int
	focused   bool
	filtering bool
	filter    textinput.Model
	viewport  viewport.Model

	// options
	width      int
	accessible bool
	theme      *Theme
	keymap     MultiSelectKeyMap
}

// NewMultiSelect returns a new multi-select field.
func NewMultiSelect[T comparable]() *MultiSelect[T] {
	filter := textinput.New()
	filter.Prompt = "/"

	return &MultiSelect[T]{
		options:   []Option[T]{},
		value:     new([]T),
		validate:  func([]T) error { return nil },
		filtering: false,
		filter:    filter,
	}
}

// Value sets the value of the multi-select field.
func (m *MultiSelect[T]) Value(value *[]T) *MultiSelect[T] {
	m.value = value
	for i, o := range m.options {
		for _, v := range *value {
			if o.Value == v {
				m.options[i].selected = true
				break
			}
		}
	}
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

	for i, o := range options {
		for _, v := range *m.value {
			if o.Value == v {
				options[i].selected = true
				break
			}
		}
	}
	m.options = options
	m.filteredOptions = options
	m.updateViewportHeight()
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

// Height sets the height of the multi-select field.
func (m *MultiSelect[T]) Height(height int) *MultiSelect[T] {
	// What we really want to do is set the height of the viewport, but we
	// need a theme applied before we can calcualate its height.
	m.height = height
	m.updateViewportHeight()
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

// Skip returns whether the multiselect should be skipped or should be blocking.
func (*MultiSelect[T]) Skip() bool {
	return false
}

// Zoom returns whether the multiselect should be zoomed.
func (*MultiSelect[T]) Zoom() bool {
	return false
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
	return []key.Binding{
		m.keymap.Toggle,
		m.keymap.Up,
		m.keymap.Down,
		m.keymap.Filter,
		m.keymap.SetFilter,
		m.keymap.ClearFilter,
		m.keymap.Prev,
		m.keymap.Submit,
		m.keymap.Next,
	}
}

// Init initializes the multi-select field.
func (m *MultiSelect[T]) Init() tea.Cmd {
	return nil
}

// Update updates the multi-select field.
func (m *MultiSelect[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Enforce height on the viewport during update as we need themes to
	// be applied before we can calculate the height.
	m.updateViewportHeight()

	var cmd tea.Cmd
	if m.filtering {
		m.filter, cmd = m.filter.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:

		m.err = nil

		switch {
		case key.Matches(msg, m.keymap.Filter):
			m.setFilter(true)
			return m, m.filter.Focus()
		case key.Matches(msg, m.keymap.SetFilter):
			if len(m.filteredOptions) <= 0 {
				m.filter.SetValue("")
				m.filteredOptions = m.options
			}
			m.setFilter(false)
		case key.Matches(msg, m.keymap.ClearFilter):
			m.filter.SetValue("")
			m.filteredOptions = m.options
			m.setFilter(false)
		case key.Matches(msg, m.keymap.Up):
			if m.filtering && msg.String() == "k" {
				break
			}

			m.cursor = max(m.cursor-1, 0)
			if m.cursor < m.viewport.YOffset {
				m.viewport.SetYOffset(m.cursor)
			}
		case key.Matches(msg, m.keymap.Down):
			if m.filtering && msg.String() == "j" {
				break
			}

			m.cursor = min(m.cursor+1, len(m.filteredOptions)-1)
			if m.cursor >= m.viewport.YOffset+m.viewport.Height {
				m.viewport.LineDown(1)
			}
		case key.Matches(msg, m.keymap.GotoTop):
			if m.filtering {
				break
			}
			m.cursor = 0
			m.viewport.GotoTop()
		case key.Matches(msg, m.keymap.GotoBottom):
			if m.filtering {
				break
			}
			m.cursor = len(m.filteredOptions) - 1
			m.viewport.GotoBottom()
		case key.Matches(msg, m.keymap.HalfPageUp):
			m.cursor = max(m.cursor-m.viewport.Height/2, 0)
			m.viewport.HalfViewUp()
		case key.Matches(msg, m.keymap.HalfPageDown):
			m.cursor = min(m.cursor+m.viewport.Height/2, len(m.filteredOptions)-1)
			m.viewport.HalfViewDown()
		case key.Matches(msg, m.keymap.Toggle):
			for i, option := range m.options {
				if option.Key == m.filteredOptions[m.cursor].Key {
					if !m.options[m.cursor].selected && m.limit > 0 && m.numSelected() >= m.limit {
						break
					}
					selected := m.options[i].selected
					m.options[i].selected = !selected
					m.filteredOptions[m.cursor].selected = !selected
				}
			}
		case key.Matches(msg, m.keymap.Prev):
			m.finalize()
			if m.err != nil {
				return m, nil
			}
			return m, PrevField
		case key.Matches(msg, m.keymap.Next, m.keymap.Submit):
			m.finalize()
			if m.err != nil {
				return m, nil
			}
			return m, NextField
		}

		if m.filtering {
			m.filteredOptions = m.options
			if m.filter.Value() != "" {
				m.filteredOptions = nil
				for _, option := range m.options {
					if m.filterFunc(option.Key) {
						m.filteredOptions = append(m.filteredOptions, option)
					}
				}
			}
			if len(m.filteredOptions) > 0 {
				m.cursor = min(m.cursor, len(m.filteredOptions)-1)
				m.viewport.SetYOffset(clamp(m.cursor, 0, len(m.filteredOptions)-m.viewport.Height))
			}
		}
	}

	return m, cmd
}

// updateViewportHeight updates the viewport size according to the Height setting
// on this multi-select field.
func (m *MultiSelect[T]) updateViewportHeight() {
	// If no height is set size the viewport to the number of options.
	if m.height <= 0 {
		m.viewport.Height = len(m.options)
		return
	}

	const minHeight = 1
	m.viewport.Height = max(minHeight, m.height-
		lipgloss.Height(m.titleView())-
		lipgloss.Height(m.descriptionView()))
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

func (m *MultiSelect[T]) activeStyles() *FieldStyles {
	theme := m.theme
	if theme == nil {
		theme = ThemeCharm()
	}
	if m.focused {
		return &theme.Focused
	}
	return &theme.Blurred
}

func (m *MultiSelect[T]) titleView() string {
	if m.title == "" {
		return ""
	}
	var (
		styles = m.activeStyles()
		sb     = strings.Builder{}
	)
	if m.filtering {
		sb.WriteString(m.filter.View())
	} else if m.filter.Value() != "" {
		sb.WriteString(styles.Title.Render(m.title) + styles.Description.Render("/"+m.filter.Value()))
	} else {
		sb.WriteString(styles.Title.Render(m.title))
	}
	if m.err != nil {
		sb.WriteString(styles.ErrorIndicator.String())
	}
	return sb.String()
}

func (m *MultiSelect[T]) descriptionView() string {
	return m.activeStyles().Description.Render(m.description)
}

func (m *MultiSelect[T]) choicesView() string {
	var (
		styles = m.activeStyles()
		c      = styles.MultiSelectSelector.String()
		sb     strings.Builder
	)
	for i, option := range m.filteredOptions {
		if m.cursor == i {
			sb.WriteString(c)
		} else {
			sb.WriteString(strings.Repeat(" ", lipgloss.Width(c)))
		}

		if m.filteredOptions[i].selected {
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

	for i := len(m.filteredOptions); i < len(m.options)-1; i++ {
		sb.WriteString("\n")
	}

	return sb.String()
}

// View renders the multi-select field.
func (m *MultiSelect[T]) View() string {
	styles := m.activeStyles()
	m.viewport.SetContent(m.choicesView())

	var sb strings.Builder
	if m.title != "" {
		sb.WriteString(m.titleView())
		sb.WriteString("\n")
	}
	if m.description != "" {
		sb.WriteString(m.descriptionView() + "\n")
	}
	sb.WriteString(m.viewport.View())
	return styles.Base.Render(sb.String())
}

func (m *MultiSelect[T]) printOptions() {
	styles := m.activeStyles()
	var sb strings.Builder

	sb.WriteString(styles.Title.Render(m.title))
	sb.WriteString("\n")

	for i, option := range m.options {
		if option.selected {
			sb.WriteString(styles.SelectedOption.Render(fmt.Sprintf("%d. %s %s", i+1, "✓", option.Key)))
		} else {
			sb.WriteString(fmt.Sprintf("%d. %s %s", i+1, " ", option.Key))
		}
		sb.WriteString("\n")
	}

	fmt.Println(sb.String())
}

// setFilter sets the filter of the select field.
func (m *MultiSelect[T]) setFilter(filter bool) {
	m.filtering = filter
	m.keymap.SetFilter.SetEnabled(filter)
	m.keymap.Filter.SetEnabled(!filter)
	m.keymap.Next.SetEnabled(!filter)
	m.keymap.Submit.SetEnabled(!filter)
	m.keymap.Prev.SetEnabled(!filter)
	m.keymap.ClearFilter.SetEnabled(!filter && m.filter.Value() != "")
}

// filterFunc returns true if the option matches the filter.
func (m *MultiSelect[T]) filterFunc(option string) bool {
	// XXX: remove diacritics or allow customization of filter function.
	return strings.Contains(strings.ToLower(option), strings.ToLower(m.filter.Value()))
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
	styles := m.activeStyles()

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

		if !m.options[choice-1].selected && m.limit > 0 && m.numSelected() >= m.limit {
			fmt.Printf("You can't select more than %d options.\n", m.limit)
			continue
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

	fmt.Println(styles.SelectedOption.Render("Selected:", strings.Join(values, ", ")+"\n"))
	return nil
}

// WithTheme sets the theme of the multi-select field.
func (m *MultiSelect[T]) WithTheme(theme *Theme) Field {
	if m.theme != nil {
		return m
	}
	m.theme = theme
	m.filter.Cursor.Style = m.theme.Focused.TextInput.Cursor
	m.filter.PromptStyle = m.theme.Focused.TextInput.Prompt
	m.updateViewportHeight()
	return m
}

// WithKeyMap sets the keymap of the multi-select field.
func (m *MultiSelect[T]) WithKeyMap(k *KeyMap) Field {
	m.keymap = k.MultiSelect
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

// WithHeight sets the height of the multi-select field.
func (m *MultiSelect[T]) WithHeight(height int) Field {
	m.height = height
	return m
}

// WithPosition sets the position of the multi-select field.
func (m *MultiSelect[T]) WithPosition(p FieldPosition) Field {
	if m.filtering {
		return m
	}
	m.keymap.Prev.SetEnabled(!p.IsFirst())
	m.keymap.Next.SetEnabled(!p.IsLast())
	m.keymap.Submit.SetEnabled(p.IsLast())
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
