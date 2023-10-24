package huh

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"
)

// Text is a form text field.
type Text struct {
	value *string
	title string

	validate func(string) error
	err      error

	textarea textarea.Model

	focused bool
	theme   *Theme
	keymap  *TextKeyMap
}

// NewText returns a new text field.
func NewText() *Text {
	text := textarea.New()
	text.ShowLineNumbers = false
	text.Prompt = ""
	text.FocusedStyle.CursorLine = lipgloss.NewStyle()

	t := &Text{
		value:    new(string),
		textarea: text,
		validate: func(string) error { return nil },
	}

	return t
}

// Value sets the value of the text field.
func (t *Text) Value(value *string) *Text {
	t.value = value
	return t
}

// Title sets the title of the text field.
func (t *Text) Title(title string) *Text {
	t.title = title
	return t
}

// CharLimit sets the character limit of the text field.
func (t *Text) CharLimit(charlimit int) *Text {
	t.textarea.CharLimit = charlimit
	return t
}

// Placeholder sets the placeholder of the text field.
func (t *Text) Placeholder(str string) *Text {
	t.textarea.Placeholder = str
	return t
}

// Validate sets the validation function of the text field.
func (t *Text) Validate(validate func(string) error) *Text {
	t.validate = validate
	return t
}

// Error returns the error of the text field.
func (t *Text) Error() error {
	return t.err
}

// Focus focuses the text field.
func (t *Text) Focus() tea.Cmd {
	t.focused = true
	return t.textarea.Focus()
}

// Blur blurs the text field.
func (t *Text) Blur() tea.Cmd {
	t.focused = false
	*t.value = t.textarea.Value()
	t.textarea.Blur()
	t.err = t.validate(*t.value)
	return nil
}

// KeyMap sets the keymap on a text field.
func (t *Text) KeyMap(k *KeyMap) Field {
	t.keymap = &k.Text
	return t
}

// KeyBinds returns the help message for the text field.
func (t *Text) KeyBinds() []key.Binding {
	return []key.Binding{t.keymap.Next, t.keymap.NewLine, t.keymap.Prev}
}

// Init initializes the text field.
func (t *Text) Init() tea.Cmd {
	t.textarea.Blur()
	return nil
}

// Update updates the text field.
func (t *Text) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	t.textarea, cmd = t.textarea.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		t.err = nil

		switch {
		case key.Matches(msg, t.keymap.Next):
			cmds = append(cmds, nextField)
		case key.Matches(msg, t.keymap.Prev):
			cmds = append(cmds, prevField)
		}
	}

	return t, tea.Batch(cmds...)
}

// View renders the text field.
func (t *Text) View() string {
	var (
		styles         FieldStyles
		textareaStyles *textarea.Style
	)
	if t.focused {
		styles = t.theme.Focused
		textareaStyles = &t.textarea.FocusedStyle
	} else {
		styles = t.theme.Blurred
		textareaStyles = &t.textarea.BlurredStyle
	}

	// NB: since the method is on a pointer receiver these are being mutated.
	// Because this runs on every render this shouldn't matter in practice,
	// however.
	textareaStyles.Placeholder = styles.TextInput.Placeholder
	textareaStyles.Text = styles.TextInput.Text
	textareaStyles.Prompt = styles.TextInput.Prompt
	textareaStyles.CursorLine = styles.TextInput.Text
	t.textarea.Cursor.Style = styles.TextInput.Cursor

	var sb strings.Builder
	sb.WriteString(styles.Title.Render(t.title))
	if t.err != nil {
		sb.WriteString(styles.ErrorIndicator.String())
	}
	sb.WriteString("\n")
	sb.WriteString(t.textarea.View())

	return styles.Base.Render(sb.String())
}

func (t *Text) Run() {
	fmt.Println(t.theme.Focused.Title.Render(t.title))
	*t.value = accessibility.PromptString("> ", t.validate)
	fmt.Println()
}

func (t *Text) Theme(theme *Theme) Field {
	t.theme = theme
	return t
}
