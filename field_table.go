package huh

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/huh/internal/accessibility" // Keep this
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	// accessor "github.com/charmbracelet/huh/accessor" removed
)

// Ensure Table implements the Field interface.
var _ Field = &Table{}

// Placeholders for functions assumed to be in the huh package or globally available
// var nextID func() int = func() int { var c int; c++; return c } // Example placeholder
// func Run(Field) error { return nil } // Example placeholder

// FieldPosition is a placeholder for field positioning data.
// This might be defined in a common package or specific to form layout.
type FieldPosition struct {
	// Example fields, replace with actual if available
	Row    int
	Column int
}

// updateTitleMsg is a message to update the table's title.
type updateTitleMsg struct{}

// updateDescriptionMsg is a message to update the table's description.
type updateDescriptionMsg struct{}

// TableKeyMap is defined in keymap.go. We use it here for the struct field.
// No, TableKeyMap is defined in keymap.go but the field `keymap` in `Table` struct is of type `TableKeyMap`
// which is correct as per previous steps. The `huh.Theme` and `huh.DefaultTheme` are from the package, not local placeholders.

// Table represents a table field that can be used to display and interact with tabular data.
type Table struct {
	id          int
	accessor    Accessor[any] // Changed from accessor.Accessor
	key         string
	title       Eval[string] // Changed from accessor.Eval
	description Eval[string] // Changed from accessor.Eval
	columns     []table.Column
	rows        []table.Row
	tableModel  table.Model
	focused     bool
	width       int
	height      int
	theme       *Theme
	keymap      TableKeyMap
	err         error
	validate    func(any) error

	// new fields
	accessible bool
	position   FieldPosition
}

// NewTable creates a new Table field.
func NewTable() *Table {
	tm := table.New() // bubbles/table model

	// Initialize keymap with default table keybindings
	defaultKeyMap := NewDefaultKeyMap() // from keymap.go

	t := &Table{
		id:          nextID(), // from huh package
		tableModel:  tm,
		title:       Eval[string]{val: ""}, // Was accessor.NewString("")
		description: Eval[string]{val: ""}, // Was accessor.NewString("")
		validate:    func(v any) error { return nil },
		keymap:      defaultKeyMap.Table,
		// theme is initially nil, will be set by activeStyles or WithTheme
	}

	// Apply initial styles based on a default theme.
	// activeStyles will initialize t.theme if it's nil.
	t.updateTableStyles()

	return t
}

// activeStyles returns the appropriate FieldStyles based on the field's focus state.
// It initializes the theme to ThemeCharm() if it's not already set.
func (t *Table) activeStyles() *FieldStyles {
	if t.theme == nil {
		t.theme = ThemeCharm() // Default theme from huh package
	}
	if t.focused {
		return &t.theme.Focused
	}
	return &t.theme.Blurred
}

// updateTableStyles applies the current theme's table styles to the underlying table.Model.
func (t *Table) updateTableStyles() {
	if t.theme == nil { // Ensure theme is initialized
		t.theme = ThemeCharm()
	}

	currentStyles := t.activeStyles()
	themedTableStyles := currentStyles.Table // These are huh.TableStyles

	// Create bubbles/table.Styles and map from our theme
	bubbleTableStyles := table.Styles{
		Header:   themedTableStyles.Header,
		Cell:     themedTableStyles.Cell,
		Selected: themedTableStyles.SelectedRow,
		// Note: bubbles/table.Styles does not have direct equivalents for
		// SelectedCell, Cursor, Border, EvenRow, OddRow from our huh.TableStyles.
		// Those would require more custom rendering if strictly needed.
		// For now, we map the primary styles.
	}
	t.tableModel.SetStyles(bubbleTableStyles)
}

// Accessor methods

// Value sets the selected row in the table.
// For now, this is a NOOP. The logic to find and select a row based on `data`
// needs to be implemented, considering the structure of `table.Row` (which is `[]string`).
func (t *Table) Value(data *any) *Table { // data is likely *table.Row or similar
	if t.accessor != nil && data != nil {
		// Attempt to set the value. If data is not of the correct type for the accessor,
		// this might be a no-op or panic depending on the accessor implementation.
		// For a table, the accessor might expect table.Row.
		// We should also update the tableModel's cursor to reflect this selection.
		if v, ok := (*data).(table.Row); ok {
			t.accessor.Set(v)
			// Find and set cursor in tableModel
			for i, r := range t.tableModel.Rows() {
				if equalRows(r, v) {
					t.tableModel.SetCursor(i)
					break
				}
			}
		} else if v, ok := (*data).([]string); ok { // Handle if *data is []string
			t.accessor.Set(table.Row(v))
			for i, r := range t.tableModel.Rows() {
				if equalRows(r, v) {
					t.tableModel.SetCursor(i)
					break
				}
			}
		}
		// If *data is nil, we might want to clear selection
	}
	return t
}

// Helper function to compare two table.Row ([]string)
func equalRows(r1, r2 table.Row) bool {
	if len(r1) != len(r2) {
		return false
	}
	for i := range r1 {
		if r1[i] != r2[i] {
			return false
		}
	}
	return true
}


// Accessor sets the accessor for the table field.
func (t *Table) Accessor(acc Accessor[any]) *Table { // Changed from accessor.Accessor
	t.accessor = acc
	return t
}

// GetValue returns the currently selected row data.
// table.Row is `[]string`.
func (t *Table) GetValue() any {
	if t.accessor != nil {
		return t.accessor.Get()
	}
	return t.tableModel.SelectedRow()
}

// GetKey returns the key of the table field.
func (t *Table) GetKey() string {
	return t.key
}

// Configuration methods

// Key sets the key of the table field.
func (t *Table) Key(key string) *Table {
	t.key = key
	return t
}

// Title sets the title of the table field.
func (t *Table) Title(title string) *Table {
	t.title = Eval[string]{val: title} // Was accessor.NewString(title)
	return t
}

// TitleFunc sets a function to dynamically update the title.
func (t *Table) TitleFunc(f func() string, bindings any) *Table {
	t.title = Eval[string]{fn: f, bindings: bindings, cache: make(map[uint64]string)} // Was accessor.NewEval(f, bindings)
	return t
}

// Description sets the description of the table field.
func (t *Table) Description(desc string) *Table {
	t.description = Eval[string]{val: desc} // Was accessor.NewString(desc)
	return t
}

// DescriptionFunc sets a function to dynamically update the description.
func (t *Table) DescriptionFunc(f func() string, bindings any) *Table {
	t.description = Eval[string]{fn: f, bindings: bindings, cache: make(map[uint64]string)} // Was accessor.NewEval(f, bindings)
	return t
}

// Columns sets the columns of the table.
func (t *Table) Columns(cols []table.Column) *Table {
	t.columns = cols
	t.tableModel.SetColumns(cols)
	return t
}

// Rows sets the rows of the table.
func (t *Table) Rows(rows []table.Row) *Table {
	t.rows = rows
	t.tableModel.SetRows(rows)
	return t
}

// Validate sets the validation function for the table field.
// The validation function receives the currently selected row (`table.Row` which is `[]string`).
func (t *Table) Validate(validate func(any) error) *Table {
	t.validate = validate
	return t
}

// Height sets the height of the table's viewport.
func (t *Table) Height(h int) *Table {
	t.height = h
	t.tableModel.SetHeight(h)
	return t
}

// Width sets the width of the table.
func (t *Table) Width(w int) *Table {
	t.width = w
	// table.Model does not have a SetWidth method directly.
	// Width is usually managed by the layout/parent components.
	// We store it for the View method or other layout calculations.
	// If direct width setting on tableModel is needed, columns widths should be adjusted.
	return t
}

// Core Interface methods

// Init initializes the table field.
// It initializes the dynamic title and description.
func (t *Table) Init() tea.Cmd {
	var cmds []tea.Cmd
	if t.title != nil {
		cmds = append(cmds, t.title.Init())
	}
	if t.description != nil {
		cmds = append(cmds, t.description.Init())
	}
	// The table model itself doesn't have an Init method that returns a command.
	// Focus is handled by the Focus method.
	return tea.Batch(cmds...)
}

// Update handles messages for the table field.
func (t *Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case updateTitleMsg:
		if t.title != nil {
			var content string
			content, cmd = t.title.Update()
			cmds = append(cmds, cmd)
			// Potentially update table title if it has one, or handle in View
		}
		return t, tea.Batch(cmds...)
	case updateDescriptionMsg:
		if t.description != nil {
			var content string
			content, cmd = t.description.Update()
			cmds = append(cmds, cmd)
			// Potentially update table description if it has one, or handle in View
		}
		return t, tea.Batch(cmds...)
	case tea.KeyMsg:
		if !t.focused {
			return t, nil
		}
		// Pass key messages to the table model for navigation
		var newTableModel table.Model
		newTableModel, cmd = t.tableModel.Update(msg)
		t.tableModel = newTableModel
		cmds = append(cmds, cmd)

		// Update accessor if a selection was made
		if t.accessor != nil {
			// Assuming SelectedRow() is the value to update
			selectedRow := t.tableModel.SelectedRow()
			if selectedRow != nil { // Check if a row is actually selected
				t.accessor.Set(selectedRow)
			}
		}

		return t, tea.Batch(cmds...)
	}

	// If the message was not a tea.KeyMsg, and not one of our update messages,
	// it might still need to be processed by the table.Model if it's a tea.WindowSizeMsg for example.
	// However, table.Model's Update typically expects tea.KeyMsg.
	// For other messages like tea.WindowSizeMsg, the parent component (group/form)
	// should handle resizing and then explicitly call t.SetWidth() or t.SetHeight().
	// For now, we only pass KeyMsg when focused.
	if t.title != nil {
		_, titleCmd := t.title.Update()
		cmds = append(cmds, titleCmd)
	}
	if t.description != nil {
		_, descCmd := t.description.Update()
		cmds = append(cmds, descCmd)
	}


	return t, tea.Batch(cmds...)
}

// View renders the table field.
func (t *Table) View() string {
	var sb strings.Builder

	currentStyles := t.activeStyles() // This ensures theme is initialized and gets focused/blurred styles

	var titleStyle, descStyle, errorStyle lipgloss.Style
	var errorIndicator string

	titleStyle = currentStyles.Title
	descStyle = currentStyles.Description
	errorStyle = currentStyles.ErrorMessage
	if t.focused {
		errorIndicator = currentStyles.ErrorIndicator.String() // Assuming ErrorIndicator is a style
	} else {
		errorIndicator = currentStyles.ErrorIndicator.String() // Or t.theme.Blurred.ErrorIndicator.String()
	}


	if t.title.String() != "" {
		sb.WriteString(titleStyle.Render(t.title.String()))
		sb.WriteString("\n")
	}
	if t.description.String() != "" {
		sb.WriteString(descStyle.Render(t.description.String()))
		sb.WriteString("\n")
	}
	// Table styles are applied to t.tableModel via updateTableStyles in Focus/Blur/WithTheme/NewTable.
	// So, t.tableModel.View() should render with the correct theme styles.

	// Apply width to the table model's base style if possible,
	// or ensure the container respects t.width.
	// For now, table.Model uses available width or its own content width.
	// If t.width is set, we could try to influence the table's container.
	// table.View() handles its own rendering based on columns, rows, styles.
	// We might need to wrap tableModel.View() in a lipgloss.Place or similar
	// if we want to strictly enforce t.width.
	// For now, let the table manage its own width.
	// Height is managed by tableModel.SetHeight().
	sb.WriteString(t.tableModel.View())

	// Error rendering part using currentStyles
	if t.err != nil {
		sb.WriteString("\n")
		// Get the string from the style for the indicator
		indicatorStr := ""
		if t.focused {
			// Render the style, which should include the string set by SetString() in the theme
			indicatorStr = currentStyles.ErrorIndicator.Render("")
		} else {
			// For blurred, it's conventional to use the blurred style's indicator
			indicatorStr = t.theme.Blurred.ErrorIndicator.Render("")
		}

		if indicatorStr != "" {
			indicatorStr += " " // Add space if indicator exists
		}
		sb.WriteString(errorStyle.Render(indicatorStr + t.err.Error()))
	}

	return sb.String()
}

// Focus focuses the table field.
func (t *Table) Focus() tea.Cmd {
	t.focused = true
	t.updateTableStyles() // Apply focused styles
	t.tableModel.Focus()  // table.Model.Focus() makes it listen to key events.
	// return accessor.EvalFocusCmd(t.title, t.description) // This needs to be huh.EvalFocusCmd or similar
	// For now, let's return nil as a placeholder if EvalFocusCmd is not readily available in package huh.
	// Or, if t.title and t.description are Eval types, they might have an Init/Focus method.
	var cmds []tea.Cmd
	if t.title.fn != nil { // Check if it's an evaluable title
		cmds = append(cmds, t.title.Init())
	}
	if t.description.fn != nil { // Check if it's an evaluable description
		cmds = append(cmds, t.description.Init())
	}
	return tea.Batch(cmds...)
}

// Blur blurs the table field.
func (t *Table) Blur() tea.Cmd {
	t.focused = false
	t.updateTableStyles() // Apply blurred styles
	t.tableModel.Blur()   // table.Model.Blur() makes it stop listening to key events.
	t.err = t.validate(t.GetValue())
	return nil
}

// Error returns the validation error of the table field.
func (t *Table) Error() error {
	return t.err
}

// Skip returns false, indicating the table field should not be skipped.
func (t *Table) Skip() bool {
	return false
}

// Zoom returns false, indicating the table field is not zoomable.
// This might change if a table needs a full-screen view.
func (t *Table) Zoom() bool {
	return false
}

// KeyBinds returns the keybindings for the table field.
// These are the bindings that are relevant for the user when the table is focused.
func (t *Table) KeyBinds() []key.Binding {
	// Note: Prev, Next, and Submit are often handled by the form/group context
	// but are included here as per the task description.
	// Depending on the desired UX, some of these might be shown conditionally
	// or not at all if they are contextually handled by the parent form.
	return []key.Binding{
		t.keymap.Up,
		t.keymap.Down,
		t.keymap.Left,
		t.keymap.Right,
		t.keymap.PageUp,
		t.keymap.PageDown,
		t.keymap.Top,
		t.keymap.Bottom,
		t.keymap.Select,
		t.keymap.Prev,
		t.keymap.Next,
		t.keymap.Submit,
	}
}

// Helper/Integration methods

// WithTheme sets the theme for the table field.
func (t *Table) WithTheme(theme *Theme) Field {
	t.theme = theme
	t.updateTableStyles()
	return t
}

// WithKeyMap sets the keymap for the table field.
// It expects a general KeyMap and will use the Table specific bindings.
func (t *Table) WithKeyMap(k *KeyMap) Field {
	if k != nil {
		t.keymap = k.Table
	}
	return t
}

// WithAccessible sets the accessible flag for the table field.
func (t *Table) WithAccessible(accessible bool) Field {
	t.accessible = accessible
	return t
}

// WithWidth sets the width of the table field.
func (t *Table) WithWidth(width int) Field {
	t.Width(width)
	return t
}

// WithHeight sets the height of the table field.
func (t *Table) WithHeight(height int) Field {
	t.Height(height)
	return t
}

// WithPosition sets the position of the table field.
func (t *Table) WithPosition(p FieldPosition) Field {
	t.position = p
	return t
}

// Run runs the table field.
// If accessible mode is enabled, it runs the accessible version of the field.
// Otherwise, it runs the bubble tea component.
func (t *Table) Run() error {
	if t.accessible {
		// Import "os" for os.Stdout, os.Stdin
		// Import "github.com/charmbracelet/huh/internal/accessibility"
		return t.runAccessible(os.Stdout, os.Stdin)
	}
	return t.run()
}

// run runs the bubble tea component for the table field.
func (t *Table) run() error {
	// This typically involves running the bubble tea program for this field.
	// In `huh`, this is often done by calling the main `Run` function for a single field.
	return Run(t) // Assumes Run(Field) is the way to run a single field.
}

// RunAccessible runs the table field in accessible mode.
// This provides a command-line interface for selecting a row.
// It is exported to allow testing from external packages.
func (t *Table) RunAccessible(w io.Writer, r io.Reader) error {
	// Imports: "fmt", "io", "strings" (strconv not used directly here)
	// "github.com/charmbracelet/huh/internal/accessibility"

	styles := t.activeStyles()
	if t.title != nil && t.title.String() != "" {
		fmt.Fprintln(w, styles.Title.Render(t.title.String()))
	}
	if t.description != nil && t.description.String() != "" {
		fmt.Fprintln(w, styles.Description.Render(t.description.String()))
	}
	fmt.Fprintln(w) // Add a blank line for spacing

	// Print headers
	if len(t.columns) > 0 {
		var headerRow strings.Builder
		for i, col := range t.columns {
			headerRow.WriteString(col.Title)
			if i < len(t.columns)-1 {
				headerRow.WriteString(" | ")
			}
		}
		fmt.Fprintln(w, styles.Table.Header.Render(headerRow.String()))
		fmt.Fprintln(w, strings.Repeat("-", headerRow.Len())) // Separator line
	}

	currentRows := t.tableModel.Rows()
	if len(currentRows) == 0 {
		fmt.Fprintln(w, "No rows to select.")
		// What should happen here? If no rows, can't select.
		// Maybe set to nil or an empty value if accessor expects it.
		if t.accessor != nil {
			// This depends on what an "empty" selection means.
			// For now, we'll assume it means no change or setting to a zero value.
			// t.accessor.Set(nil) // Or an empty table.Row / []string
		}
		return nil
	}

	for i, row := range currentRows {
		var rowStr strings.Builder
		for j, cell := range row {
			rowStr.WriteString(cell)
			if j < len(row)-1 {
				rowStr.WriteString(" | ")
			}
		}
		// Apply Cell style if needed, but for accessibility, raw text is often better.
		// For now, just print the text.
		fmt.Fprintf(w, "%d. %s\n", i+1, rowStr.String())
	}
	fmt.Fprintln(w)

	defaultChoice := 0 // 0 means no default or first item if 1-indexed
	currentVal := t.accessor.Get()
	if currentVal != nil {
		if selectedRow, ok := currentVal.(table.Row); ok {
			for i, r := range currentRows {
				// Simple comparison for table.Row (which is []string)
				if len(selectedRow) == len(r) {
					match := true
					for k := range selectedRow {
						if selectedRow[k] != r[k] {
							match = false
							break
						}
					}
					if match {
						defaultChoice = i + 1 // 1-indexed
						break
					}
				}
			}
		}
	}

	prompt := fmt.Sprintf("Choose a row (1-%d)", len(currentRows))
	if defaultChoice > 0 {
		prompt = fmt.Sprintf("%s [%d]", prompt, defaultChoice)
	}

	for {
		choice, err := accessibility.PromptInt(prompt+": ", r, w, defaultChoice, 1, len(currentRows))
		if err != nil {
			// Handle EOF or other read errors
			// accessibility.PromptInt should ideally handle basic re-prompting on parse error.
			// If it returns an error, it's likely a more serious issue (e.g., EOF).
			fmt.Fprintln(w, styles.ErrorMessage.Render("Error reading input: "+err.Error()))
			return err
		}

		selectedIndex := choice - 1 // Convert 1-indexed to 0-indexed
		selectedRowData := currentRows[selectedIndex]

		// Perform field validation
		validationErr := t.validate(selectedRowData)
		if validationErr != nil {
			errorMsg := styles.ErrorMessage.Render(validationErr.Error())
			indicator := styles.ErrorIndicator.Render("")
			if indicator != "" {
				indicator += " "
			}
			fmt.Fprintln(w, indicator+errorMsg)
			// Re-prompt by continuing the loop
			defaultChoice = choice // Keep the last tried choice as default for next prompt
			continue
		}

		// If validation passes, set the value
		if t.accessor != nil {
			t.accessor.Set(selectedRowData)
		}
		break // Exit loop on successful selection and validation
	}

	return nil
}
