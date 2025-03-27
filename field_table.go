package huh

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"
)

// Column defines the table structure and a getter
// function allowing to retrieve the cell value from the
// row structure.
type Column[T any] struct {
	table.Column
	Get func(T) any
}

// NewColumn create a new Column for the table field.
func NewColumn[T any](title string, width int, get func(T) any) Column[T] {
	return Column[T]{
		Column: table.Column{
			Title: title,
			Width: width,
		},
		Get: get,
	}
}

// Table is a select field rendered as a table.
type Table[T any, K comparable] struct {
	id       int
	accessor Accessor[T]
	key      string

	table table.Model

	title       Eval[string]
	description Eval[string]
	columns     []Column[T]
	options     Eval[[]TableOption[T, K]]

	validate func(T) error
	err      error

	height     int
	selected   int
	accessible bool
	spinner    spinner.Model

	theme  *Theme
	keymap TableKeyMap
}

// NewTable creates a new table field.
func NewTable[T any, K comparable]() *Table[T, K] {
	t := table.New()

	s := spinner.New(spinner.WithSpinner(spinner.Line))

	return &Table[T, K]{
		table:       t,
		accessor:    &EmbeddedAccessor[T]{},
		validate:    func(T) error { return nil },
		options:     Eval[[]TableOption[T, K]]{cache: make(map[uint64][]TableOption[T, K])},
		title:       Eval[string]{cache: make(map[uint64]string)},
		description: Eval[string]{cache: make(map[uint64]string)},
		spinner:     s,
	}
}

// Value sets the value of the table field.
func (t *Table[T, K]) Value(value *T) *Table[T, K] {
	return t.Accessor(NewPointerAccessor(value))
}

// Accessor sets the accessor of the table field.
func (t *Table[T, K]) Accessor(accessor Accessor[T]) *Table[T, K] {
	t.accessor = accessor
	t.selectValue(t.accessor.Get())
	t.updateValue()
	return t
}

func (t *Table[T, K]) selectValue(value T) {
	for i, o := range t.options.val {
		if o.Key() == o.key(value) {
			t.selected = i
			break
		}
	}
}

// Key sets the key of the table field which can be used to retrieve the value
// after submission.
func (t *Table[T, K]) Key(key string) *Table[T, K] {
	t.key = key
	return t
}

// Title sets the title of the table field.
//
// This title will be static, for dynamic titles use `TitleFunc`.
func (t *Table[T, K]) Title(title string) *Table[T, K] {
	t.title.val = title
	t.title.fn = nil
	return t
}

// TitleFunc sets the title func of the table field.
//
// This TitleFunc will be re-evaluated when the binding of the TitleFunc
// changes. This when you want to display dynamic content and update the title
// when another part of your form changes.
//
// See README#Dynamic for more usage information.
func (t *Table[T, K]) TitleFunc(f func() string, bindings any) *Table[T, K] {
	t.title.fn = f
	t.title.bindings = bindings
	return t
}

// Description sets the description of the table field.
//
// This description will be static, for dynamic descriptions use `DescriptionFunc`.
func (t *Table[T, K]) Description(description string) *Table[T, K] {
	t.description.val = description
	return t
}

// DescriptionFunc sets the description func of the table field.
//
// This DescriptionFunc will be re-evaluated when the binding of the
// DescriptionFunc changes. This is useful when you want to display dynamic
// content and update the description when another part of your form changes.
//
// See README#Dynamic for more usage information.
func (t *Table[T, K]) DescriptionFunc(f func() string, bindings any) *Table[T, K] {
	t.description.fn = f
	t.description.bindings = bindings
	return t
}

// Columns set the table columns.
func (t *Table[T, K]) Columns(columns ...Column[T]) *Table[T, K] {
	t.columns = columns
	tableColumns := make([]table.Column, 0, len(t.columns))
	for _, c := range t.columns {
		tableColumns = append(tableColumns, c.Column)
	}
	t.table.SetColumns(tableColumns)
	return t
}

// Options sets the options of the table field.
func (t *Table[T, K]) Options(options ...TableOption[T, K]) *Table[T, K] {
	if len(options) <= 0 {
		return t
	}
	t.options.val = options

	// Set the cursor to the existing value or the last selected option.
	for i, option := range options {
		if option.Key() == option.key(t.accessor.Get()) {
			t.selected = i
			break
		} else if option.selected {
			t.selected = i
		}
	}

	rows := make([]table.Row, 0, len(options))
	for _, option := range options {
		row := make(table.Row, 0, len(t.columns))
		for _, c := range t.columns {
			row = append(row, fmt.Sprintf("%v", c.Get(option.Value)))
		}
		rows = append(rows, row)
	}
	t.table.SetRows(rows)
	t.updateValue()

	return t
}

// OptionsFunc sets the options func of the table field.
//
// This OptionsFunc will be re-evaluated when the binding of the OptionsFunc
// changes. This is useful when you want to display dynamic content and update
// the options when another part of your form changes.
func (t *Table[T, K]) OptionsFunc(f func() []TableOption[T, K], bindings any) *Table[T, K] {
	t.options.fn = f
	t.options.bindings = bindings
	// If there is no height set, we should attach a static height since these
	// options are possibly dynamic.
	if t.height <= 0 {
		t.height = defaultHeight
	}
	return t
}

// Height sets the height of the table field. If the number of options exceeds
// the height, the table field will become scrollable.
func (t *Table[T, K]) Height(height int) *Table[T, K] {
	t.height = height
	t.table.SetHeight(t.tableHeight())
	return t
}

// Width sets the width of the table field.
func (t *Table[T, K]) Width(width int) *Table[T, K] {
	t.table.SetWidth(width)
	return t
}

// Validate sets the validation function of the select field.
func (t *Table[T, K]) Validate(validate func(T) error) *Table[T, K] {
	t.validate = validate
	return t
}

// Error returns the error of the table field.
func (t *Table[T, K]) Error() error { return t.err }

// Skip returns whether the select should be skipped or should be blocking.
func (*Table[T, K]) Skip() bool { return false }

// Zoom returns whether the input should be zoomed.
func (*Table[T, K]) Zoom() bool { return false }

// Focus focuses the table field.
func (t *Table[T, K]) Focus() tea.Cmd {
	t.table.Focus()
	return nil
}

// Blur blurs the table field.
func (t *Table[T, K]) Blur() tea.Cmd {
	value := t.accessor.Get()
	t.table.Blur()
	t.err = t.validate(value)
	return nil
}

// KeyBinds returns the help keybindings for the table field.
func (t *Table[T, K]) KeyBinds() []key.Binding {
	return []key.Binding{
		t.keymap.LineUp,
		t.keymap.LineDown,
		t.keymap.GotoBottom,
		t.keymap.GotoTop,
		t.keymap.HalfPageDown,
		t.keymap.HalfPageUp,
		t.keymap.PageDown,
		t.keymap.PageUp,
		t.keymap.Prev,
		t.keymap.Next,
		t.keymap.Submit,
	}
}

// Init initializes the table field.
func (t *Table[T, K]) Init() tea.Cmd {
	return nil
}

// Update updates the table field.
func (t *Table[T, K]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case updateFieldMsg:
		if ok, hash := t.title.shouldUpdate(); ok {
			t.title.bindingsHash = hash
			if !t.title.loadFromCache() {
				t.title.loading = true
				cmds = append(cmds, func() tea.Msg {
					return updateTitleMsg{id: t.id, title: t.title.fn(), hash: hash}
				})
			}
		}
		if ok, hash := t.description.shouldUpdate(); ok {
			t.description.bindingsHash = hash
			if !t.description.loadFromCache() {
				t.description.loading = true
				cmds = append(cmds, func() tea.Msg {
					return updateDescriptionMsg{id: t.id, description: t.description.fn(), hash: hash}
				})
			}
		}
		if ok, hash := t.options.shouldUpdate(); ok {
			t.options.bindingsHash = hash
			if t.options.loadFromCache() {
				t.selected = clamp(t.selected, 0, len(t.options.val)-1)
			} else {
				t.options.loading = true
				t.options.loadingStart = time.Now()
				cmds = append(cmds, func() tea.Msg {
					return updateTableOptionsMsg[T, K]{id: t.id, hash: hash, options: t.options.fn()}
				}, t.spinner.Tick)
			}
		}
		return t, tea.Batch(cmds...)
	case spinner.TickMsg:
		if !t.options.loading {
			break
		}
		t.spinner, cmd = t.spinner.Update(msg)
		return t, cmd

	case updateTitleMsg:
		if msg.id == t.id && msg.hash == t.title.bindingsHash {
			t.title.update(msg.title)
		}
	case updateDescriptionMsg:
		if msg.id == t.id && msg.hash == t.description.bindingsHash {
			t.description.update(msg.description)
		}
	case updateTableOptionsMsg[T, K]:
		if msg.id == t.id && msg.hash == t.options.bindingsHash {
			t.options.update(msg.options)

			// since we're updating the options, we need to update the selected cursor
			// position and filteredOptions.
			t.selected = clamp(t.selected, 0, len(msg.options)-1)
			t.updateValue()
		}
	case tea.KeyMsg:
		t.err = nil
		switch {
		case key.Matches(msg, t.keymap.Prev):
			if t.selected >= len(t.options.val) {
				break
			}
			t.updateValue()
			t.err = t.validate(t.accessor.Get())
			if t.err != nil {
				return t, nil
			}
			t.updateValue()
			return t, PrevField
		case key.Matches(msg, t.keymap.Next, t.keymap.Submit):
			if t.selected >= len(t.options.val) {
				break
			}
			t.updateValue()
			t.err = t.validate(t.accessor.Get())
			if t.err != nil {
				return t, nil
			}
			t.updateValue()
			return t, NextField
		}
	}
	t.table, cmd = t.table.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	t.selected = t.table.Cursor()
	t.updateValue()
	return t, tea.Batch(cmds...)
}

func (t *Table[T, K]) updateValue() {
	if t.selected < len(t.options.val) && t.selected >= 0 {
		t.accessor.Set(t.options.val[t.selected].Value)
	}
}

func (t *Table[T, K]) tableHeight() int {
	height := t.height
	if t.title.val != "" {
		height -= lipgloss.Height(t.title.val)
	}
	if t.description.val != "" {
		height -= lipgloss.Height(t.description.val)
	}
	if t.title.val != "" || t.description.val != "" {
		height-- // Space between table and title/description
	}
	return height
}

func (t *Table[T, K]) activeStyles() *FieldStyles {
	theme := t.theme
	if theme == nil {
		theme = ThemeCharm()
	}
	if t.table.Focused() {
		return &theme.Focused
	}
	return &theme.Blurred
}

func (t *Table[T, K]) titleView() string {
	var (
		styles = t.activeStyles()
		sb     = strings.Builder{}
	)
	sb.WriteString(styles.Title.Render(t.title.val))
	if t.err != nil {
		sb.WriteString(styles.ErrorIndicator.String())
	}
	return sb.String()
}

func (t *Table[T, K]) descriptionView() string {
	return t.activeStyles().Description.Render(t.description.val)
}

// View renders the table field.
func (t *Table[T, K]) View() string {
	styles := t.activeStyles()

	var sb strings.Builder
	if t.title.val != "" || t.title.fn != nil {
		sb.WriteString(t.titleView())
		sb.WriteString("\n")
	}
	if t.description.val != "" || t.description.fn != nil {
		sb.WriteString(t.descriptionView())
		sb.WriteString("\n")
	}
	if t.title.val != "" || t.description.val != "" {
		sb.WriteString("\n")
	}
	t.table.SetStyles(styles.Table)
	t.table.SetHeight(t.tableHeight())
	sb.WriteString(t.table.View())
	return styles.Base.Render(sb.String())
}

// Run runs the table field.
func (t *Table[T, K]) Run() error {
	if t.accessible {
		return t.runAccessible()
	}
	return Run(t)
}

// runAccessible runs an accessible table field.
func (t *Table[T, K]) runAccessible() error {
	var sb strings.Builder
	styles := t.activeStyles()
	if t.title.val != "" {
		sb.WriteString(styles.Title.Render(t.title.val) + "\n")
	}
	if t.description.val != "" {
		sb.WriteString(styles.Description.Render(t.description.val) + "\n")
	}
	if t.title.val != "" || t.description.val != "" {
		sb.WriteString("\n")
	}

	bold := lipgloss.NewStyle().Bold(true)

	// Header
	sb.WriteString("    # ")
	for _, col := range t.columns {
		format := fmt.Sprintf("%%-%ds ", col.Width)
		sb.WriteString(bold.Render(fmt.Sprintf(format, col.Title)))
	}
	sb.WriteString("\n")

	// Rows
	for i, option := range t.options.val {
		sb.WriteString(fmt.Sprintf("%5d ", i+1))
		for _, col := range t.columns {
			format := fmt.Sprintf("%%-%dv ", col.Width)
			sb.WriteString(fmt.Sprintf(format, col.Get(option.Value)))
		}
		sb.WriteString("\n")
	}

	fmt.Println(sb.String())

	for {
		choice := accessibility.PromptInt("Choose: ", 1, len(t.options.val))
		option := t.options.val[choice-1]
		if err := t.validate(option.Value); err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(styles.SelectedOption.Render(fmt.Sprintf("Chose: %v\n", option.Key())))
		t.accessor.Set(option.Value)
		break
	}

	return nil
}

// WithTheme sets the theme of the select field.
func (t *Table[T, K]) WithTheme(theme *Theme) Field {
	if t.theme != nil {
		return t
	}
	t.theme = theme
	t.table.UpdateViewport()
	return t
}

// WithKeyMap sets the keymap on a table field.
func (t *Table[T, K]) WithKeyMap(k *KeyMap) Field {
	t.keymap = k.Table
	return t
}

// WithAccessible sets the accessible mode of the select field.
func (t *Table[T, K]) WithAccessible(accessible bool) Field {
	t.accessible = accessible
	return t
}

// WithWidth sets the width of the table field.
func (t *Table[T, K]) WithWidth(width int) Field {
	t.table.SetWidth(width)
	return t
}

// WithHeight sets the height of the table field.
func (t *Table[T, K]) WithHeight(height int) Field {
	return t.Height(height)
}

// WithPosition sets the position of the table field.
func (t *Table[T, K]) WithPosition(p FieldPosition) Field {
	t.keymap.Prev.SetEnabled(!p.IsFirst())
	t.keymap.Next.SetEnabled(!p.IsLast())
	t.keymap.Submit.SetEnabled(p.IsLast())
	return t
}

// GetKey returns the key of the field.
func (t *Table[T, K]) GetKey() string { return t.key }

// GetValue returns the value of the field.
func (t *Table[T, K]) GetValue() any {
	return t.accessor.Get()
}
