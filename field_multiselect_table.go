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
	"github.com/mattn/go-runewidth"
)

// Row is an option for select fields.
type Row struct {
	Key      string
	Values   []string
	selected bool
}

type StringSlice []string

type FilterFunc func(filter string, option Row) bool

// MultiSelectTable is a form multi-select field.
type MultiSelectTable struct {
	value *[]string
	key   string

	// customization
	title           string
	description     string
	options         []Row
	filterable      bool
	filteredOptions []Row
	filterFunc      FilterFunc
	limit           int
	height          int
	cols            []Column

	// error handling
	validate func([]string) error
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
	keymap     *MultiSelectKeyMap
}

// Column defines the table structure.
type Column struct {
	Title string
	Width int
}

// NewMultiSelectTable returns a new multi-select field.
func NewMultiSelectTable() *MultiSelectTable {
	filter := textinput.New()
	filter.Prompt = "/"

	return &MultiSelectTable{
		options:   []Row{},
		value:     new([]string),
		validate:  func([]string) error { return nil },
		filtering: false,
		filter:    filter,
	}
}

// Value sets the value of the multi-select field.
func (m *MultiSelectTable) Value(value *[]string) *MultiSelectTable {
	m.value = value
	for i, o := range m.options {
		for _, v := range *value {
			if o.Key == v {
				m.options[i].selected = true
				break
			}
		}
	}
	return m
}

// Key sets the key of the select field which can be used to retrieve the value
// after submission.
func (m *MultiSelectTable) Key(key string) *MultiSelectTable {
	m.key = key
	return m
}

// Title sets the title of the multi-select field.
func (m *MultiSelectTable) Title(title string) *MultiSelectTable {
	m.title = title
	return m
}

// Description sets the description of the multi-select field.
func (m *MultiSelectTable) Description(description string) *MultiSelectTable {
	m.description = description
	return m
}

// WithColumns sets the table columns (headers).
func (m *MultiSelectTable) Columns(cols []Column) *MultiSelectTable {
	if len(cols) <= 0 {
		return m
	}
	m.cols = cols
	m.updateViewportHeight()
	return m
}

// Options sets the options of the multi-select field.
func (m *MultiSelectTable) Options(options ...Row) *MultiSelectTable {
	if len(options) <= 0 {
		return m
	}

	for i, o := range options {
		for _, v := range *m.value {
			if o.Key == v {
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
func (m *MultiSelectTable) Filterable(filterable bool) *MultiSelectTable {
	m.filterable = filterable
	return m
}

// Limit sets the limit of the multi-select field.
func (m *MultiSelectTable) Limit(limit int) *MultiSelectTable {
	m.limit = limit
	return m
}

// Height sets the height of the multi-select field.
func (m *MultiSelectTable) Height(height int) *MultiSelectTable {
	// What we really want to do is set the height of the viewport, but we
	// need a theme applied before we can calcualate its height.
	m.height = height
	m.updateViewportHeight()
	return m
}

// Validate sets the validation function of the multi-select field.
func (m *MultiSelectTable) Validate(validate func([]string) error) *MultiSelectTable {
	m.validate = validate
	return m
}

// Error returns the error of the multi-select field.
func (m *MultiSelectTable) Error() error {
	return m.err
}

// Focus focuses the multi-select field.
func (m *MultiSelectTable) Focus() tea.Cmd {
	m.focused = true
	return nil
}

// Blur blurs the multi-select field.
func (m *MultiSelectTable) Blur() tea.Cmd {
	m.focused = false
	return nil
}

// KeyBinds returns the help message for the multi-select field.
func (m *MultiSelectTable) KeyBinds() []key.Binding {
	return []key.Binding{m.keymap.Toggle, m.keymap.Up, m.keymap.Down, m.keymap.Filter, m.keymap.SetFilter, m.keymap.ClearFilter, m.keymap.Next, m.keymap.Prev}
}

// Init initializes the multi-select field.
func (m *MultiSelectTable) Init() tea.Cmd {
	return nil
}

// Update updates the multi-select field.
func (m *MultiSelectTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, prevField
		case key.Matches(msg, m.keymap.Next):
			m.finalize()
			if m.err != nil {
				return m, nil
			}
			return m, nextField
		}

		if m.filtering {
			m.filteredOptions = m.options
			if m.filter.Value() != "" {
				m.filteredOptions = nil
				for _, option := range m.options {
					if m.filterRow(option) {
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
func (m *MultiSelectTable) updateViewportHeight() {
	// min accommodates the

	// If no height is set size the viewport to the number of options plus the min height.
	if m.height <= 0 {
		m.viewport.Height = len(m.options)
		return
	}

	// Wait until the theme has appied or things'll panic.
	if m.theme == nil {
		return
	}
	const minHeight = 1
	m.viewport.Height = max(minHeight, m.height-
		lipgloss.Height(m.titleView())-
		lipgloss.Height(m.descriptionView())-
		lipgloss.Height(m.headersView()))
}

func (m *MultiSelectTable) numSelected() int {
	var count int
	for _, o := range m.options {
		if o.selected {
			count++
		}
	}
	return count
}

func (m *MultiSelectTable) finalize() {
	*m.value = make([]string, 0)
	for _, option := range m.options {
		if option.selected {
			*m.value = append(*m.value, option.Key)
		}
	}
	m.err = m.validate(*m.value)
}

func (m *MultiSelectTable) activeStyles() *FieldStyles {
	if m.focused {
		return &m.theme.Focused
	}
	return &m.theme.Blurred
}

func (m *MultiSelectTable) titleString() string {
	return m.title + fmt.Sprintf(" (selected: %v)", m.numSelected())
}
func (m *MultiSelectTable) titleView() string {
	var (
		styles = m.activeStyles()
		sb     = strings.Builder{}
	)
	if m.filtering {
		sb.WriteString(m.filter.View())
	} else if m.filter.Value() != "" {
		sb.WriteString(styles.Title.Render(m.titleString()) + styles.Description.Render("/"+m.filter.Value()))
	} else {
		sb.WriteString(styles.Title.Render(m.titleString()))
	}
	if m.err != nil {
		sb.WriteString(styles.ErrorIndicator.String())
	}
	return sb.String()
}

func (m *MultiSelectTable) descriptionView() string {
	return m.activeStyles().Description.Render(m.description)
}

func (m *MultiSelectTable) renderRow(row Row) string {

	var s = make([]string, 0, len(m.cols))
	for i, value := range row.Values {
		style := lipgloss.NewStyle().Width(m.cols[i].Width).MaxWidth(m.cols[i].Width).Inline(true)
		renderedCell := m.activeStyles().Option.Render(style.Render(runewidth.Truncate(value, m.cols[i].Width, "…")))
		s = append(s, renderedCell)
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, s...)
}

func (m *MultiSelectTable) choicesView() string {
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
			sb.WriteString(styles.SelectedOption.Render(m.renderRow(option)))
		} else {
			sb.WriteString(styles.UnselectedPrefix.String())
			sb.WriteString(styles.UnselectedOption.Render(m.renderRow(option)))
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

func (m *MultiSelectTable) headersView() string {
	styles := m.activeStyles()
	var s = make([]string, 0, len(m.cols))
	for _, col := range m.cols {
		style := lipgloss.NewStyle().Width(col.Width).MaxWidth(col.Width).Inline(true)
		renderedCell := style.Render(runewidth.Truncate(col.Title, col.Width, "…"))
		s = append(s, styles.Header.Render(renderedCell))
	}

	prefix := strings.Repeat(" ", lipgloss.Width(styles.MultiSelectSelector.String())+lipgloss.Width(styles.SelectedPrefix.String()))

	return prefix + lipgloss.JoinHorizontal(lipgloss.Left, s...)
}

// View renders the multi-select field.
func (m *MultiSelectTable) View() string {
	styles := m.activeStyles()
	m.viewport.SetContent(m.choicesView())

	var sb strings.Builder

	sb.WriteString(m.titleView())
	if m.err != nil {
		sb.WriteString(styles.ErrorIndicator.String())
	}
	sb.WriteString("\n")
	if m.description != "" {
		sb.WriteString(m.descriptionView() + "\n\n")
	}

	sb.WriteString(m.headersView() + "\n")
	sb.WriteString(m.viewport.View())
	return styles.Base.Render(sb.String())
}

func (m *MultiSelectTable) printOptions() {
	var (
		sb strings.Builder
	)

	sb.WriteString(m.theme.Focused.Title.Render(m.titleString()))
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

// setFilter sets the filter of the select field.
func (m *MultiSelectTable) setFilter(filter bool) {
	m.filtering = filter
	m.keymap.SetFilter.SetEnabled(filter)
	m.keymap.Filter.SetEnabled(!filter)
	m.keymap.ClearFilter.SetEnabled(!filter && m.filter.Value() != "")
}

// filterFunc returns true if the option matches the filter.
func (m *MultiSelectTable) filterRow(option Row) bool {
	if m.filterFunc != nil {
		return m.filterFunc(m.filter.Value(), option)
	}
	// XXX: remove diacritics or allow customization of filter function.
	return strings.Contains(strings.ToLower(option.Key), strings.ToLower(m.filter.Value()))
}

// Run runs the multi-select field.
func (m *MultiSelectTable) Run() error {
	if m.accessible {
		return m.runAccessible()
	}
	return Run(m)
}

// runAccessible() runs the multi-select field in accessible mode.
func (m *MultiSelectTable) runAccessible() error {
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
			*m.value = append(*m.value, option.Key)
			values = append(values, option.Key)
		}
	}

	fmt.Println(m.theme.Focused.SelectedOption.Render("Selected:", strings.Join(values, ", ")+"\n"))
	return nil
}

// WithTheme sets the theme of the multi-select field.
func (m *MultiSelectTable) WithTheme(theme *Theme) Field {
	m.theme = theme
	m.filter.Cursor.Style = m.theme.Focused.TextInput.Cursor
	m.filter.PromptStyle = m.theme.Focused.TextInput.Prompt
	m.updateViewportHeight()
	return m
}

// WithKeyMap sets the keymap of the multi-select field.
func (m *MultiSelectTable) WithKeyMap(k *KeyMap) Field {
	m.keymap = &k.MultiSelect
	return m
}

// WithFilterFunc sets the filterFunc of the multi-select field.
func (m *MultiSelectTable) WithFilterFunc(f FilterFunc) Field {
	m.filterFunc = f
	return m
}

// WithAccessible sets the accessible mode of the multi-select field.
func (m *MultiSelectTable) WithAccessible(accessible bool) Field {
	m.accessible = accessible
	return m
}

// WithWidth sets the width of the multi-select field.
func (m *MultiSelectTable) WithWidth(width int) Field {
	m.width = width
	return m
}

// GetKey returns the multi-select's key.
func (m *MultiSelectTable) GetKey() string {
	return m.key
}

// GetValue returns the multi-select's value.
func (m *MultiSelectTable) GetValue() any {
	return *m.value
}
