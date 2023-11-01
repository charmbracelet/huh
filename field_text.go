package huh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"
)

// Text is a form text field. It allows for a multi-line string input.
type Text struct {
	value *string

	// error handling
	validate func(string) error
	err      error

	// model
	textarea textarea.Model

	// customization
	title           string
	description     string
	editorCmd       string
	editorArgs      []string
	editorExtension string

	// state
	focused bool

	// form options
	width      int
	accessible bool
	theme      *Theme
	keymap     *TextKeyMap
}

// NewText returns a new text field.
func NewText() *Text {
	text := textarea.New()
	text.ShowLineNumbers = false
	text.Prompt = ""
	text.FocusedStyle.CursorLine = lipgloss.NewStyle()

	editorCmd, editorArgs := getEditor()

	t := &Text{
		value:           new(string),
		textarea:        text,
		validate:        func(string) error { return nil },
		editorCmd:       editorCmd,
		editorArgs:      editorArgs,
		editorExtension: "md",
	}

	return t
}

// Value sets the value of the text field.
func (t *Text) Value(value *string) *Text {
	t.value = value
	t.textarea.SetValue(*value)
	return t
}

// Title sets the title of the text field.
func (t *Text) Title(title string) *Text {
	t.title = title
	return t
}

// Description sets the description of the text field.
func (t *Text) Description(description string) *Text {
	t.description = description
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

const defaultEditor = "nano"

// getEditor returns the editor command and arguments.
func getEditor() (string, []string) {
	editor := strings.Fields(os.Getenv("EDITOR"))
	if len(editor) > 0 {
		return editor[0], editor[1:]
	}
	return defaultEditor, nil
}

// Editor specifies which editor to use.
//
// The first argument provided is used as the editor command (vim, nvim, nano, etc...)
// The following (optional) arguments provided are passed as arguments to the editor command.
//
//   .Editor("vim")
//   .Editor("vim", "+10")
//
func (t *Text) Editor(editor ...string) *Text {
	if len(editor) > 0 {
		t.editorCmd = editor[0]
	}
	if len(editor) > 1 {
		t.editorArgs = editor[1:]
	}
	return t
}

// EditorExtension specifies arguments to pass into the editor.
func (t *Text) EditorExtension(extension string) *Text {
	t.editorExtension = extension
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

// KeyBinds returns the help message for the text field.
func (t *Text) KeyBinds() []key.Binding {
	return []key.Binding{t.keymap.Next, t.keymap.Editor, t.keymap.Prev}
}

type updateValueMsg []byte

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
	case updateValueMsg:
		t.textarea.SetValue(string(msg))
	case tea.KeyMsg:
		t.err = nil

		switch {
		case key.Matches(msg, t.keymap.Editor):
			ext := strings.TrimPrefix(t.editorExtension, ".")
			tmpFile, _ := os.CreateTemp(os.TempDir(), "*."+ext)
			cmd := exec.Command(t.editorCmd, append(t.editorArgs, tmpFile.Name())...)
			_ = os.WriteFile(tmpFile.Name(), []byte(t.textarea.Value()), 0644)
			cmds = append(cmds, tea.ExecProcess(cmd, func(error) tea.Msg {
				content, _ := os.ReadFile(tmpFile.Name())
				return updateValueMsg(content)
			}))
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
	if t.title != "" {
		sb.WriteString(styles.Title.Render(t.title))
		if t.err != nil {
			sb.WriteString(styles.ErrorIndicator.String())
		}
		sb.WriteString("\n")
	}
	if t.description != "" {
		sb.WriteString(styles.Description.Render(t.description))
		sb.WriteString("\n")
	}
	sb.WriteString(t.textarea.View())

	return styles.Base.Render(sb.String())
}

// Run runs the text field.
func (t *Text) Run() error {
	if t.accessible {
		return t.runAccessible()
	}
	return Run(t)
}

// runAccessible runs an accessible text field.
func (t *Text) runAccessible() error {
	fmt.Println(t.theme.Focused.Title.Render(t.title))
	*t.value = accessibility.PromptString("> ", t.validate)
	fmt.Println()
	return nil
}

// WithTheme sets the theme on a text field.
func (t *Text) WithTheme(theme *Theme) Field {
	t.theme = theme
	return t
}

// WithKeyMap sets the keymap on a text field.
func (t *Text) WithKeyMap(k *KeyMap) Field {
	t.keymap = &k.Text
	return t
}

// WithAccessible sets the accessible mode of the text field.
func (t *Text) WithAccessible(accessible bool) Field {
	t.accessible = accessible
	return t
}

// WithWidth sets the width of the text field.
func (t *Text) WithWidth(width int) Field {
	t.width = width
	return t
}
