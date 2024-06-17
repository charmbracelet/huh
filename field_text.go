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
	key   string

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
	keymap     TextKeyMap
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

// Key sets the key of the text field.
func (t *Text) Key(key string) *Text {
	t.key = key
	return t
}

// Title sets the title of the text field.
func (t *Text) Title(title string) *Text {
	t.title = title
	return t
}

// Lines sets the number of lines to show of the text field.
func (t *Text) Lines(lines int) *Text {
	t.textarea.SetHeight(lines)
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

// ShowLineNumbers sets whether or not to show line numbers.
func (t *Text) ShowLineNumbers(show bool) *Text {
	t.textarea.ShowLineNumbers = show
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

// Skip returns whether the textarea should be skipped or should be blocking.
func (*Text) Skip() bool {
	return false
}

// Zoom returns whether the note should be zoomed.
func (*Text) Zoom() bool {
	return false
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
	return []key.Binding{t.keymap.NewLine, t.keymap.Editor, t.keymap.Prev, t.keymap.Submit, t.keymap.Next}
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
	*t.value = t.textarea.Value()

	switch msg := msg.(type) {
	case updateValueMsg:
		t.textarea.SetValue(string(msg))
		t.textarea, cmd = t.textarea.Update(msg)
		cmds = append(cmds, cmd)
		*t.value = t.textarea.Value()

	case tea.KeyMsg:
		t.err = nil

		switch {
		case key.Matches(msg, t.keymap.Editor):
			ext := strings.TrimPrefix(t.editorExtension, ".")
			tmpFile, _ := os.CreateTemp(os.TempDir(), "*."+ext)
			cmd := exec.Command(t.editorCmd, append(t.editorArgs, tmpFile.Name())...)
			_ = os.WriteFile(tmpFile.Name(), []byte(t.textarea.Value()), 0600)
			cmds = append(cmds, tea.ExecProcess(cmd, func(error) tea.Msg {
				content, _ := os.ReadFile(tmpFile.Name())
				_ = os.Remove(tmpFile.Name())
				return updateValueMsg(content)
			}))
		case key.Matches(msg, t.keymap.Next, t.keymap.Submit):
			value := t.textarea.Value()
			t.err = t.validate(value)
			if t.err != nil {
				return t, nil
			}
			cmds = append(cmds, NextField)
		case key.Matches(msg, t.keymap.Prev):
			value := t.textarea.Value()
			t.err = t.validate(value)
			if t.err != nil {
				return t, nil
			}
			cmds = append(cmds, PrevField)
		}
	}

	return t, tea.Batch(cmds...)
}

func (t *Text) activeStyles() *FieldStyles {
	theme := t.theme
	if theme == nil {
		theme = ThemeCharm()
	}
	if t.focused {
		return &theme.Focused
	}
	return &theme.Blurred
}

func (t *Text) activeTextAreaStyles() *textarea.Style {
	if t.theme == nil {
		return &t.textarea.BlurredStyle
	}
	if t.focused {
		return &t.textarea.FocusedStyle
	}
	return &t.textarea.BlurredStyle
}

// View renders the text field.
func (t *Text) View() string {
	var styles = t.activeStyles()
	var textareaStyles = t.activeTextAreaStyles()

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
	styles := t.activeStyles()
	fmt.Println(styles.Title.Render(t.title))
	fmt.Println()
	*t.value = accessibility.PromptString("Input: ", func(input string) error {
		if err := t.validate(input); err != nil {
			// Handle the error from t.validate, return it
			return err
		}

		if len(input) > t.textarea.CharLimit {
			return fmt.Errorf("Input cannot exceed %d characters", t.textarea.CharLimit)
		}
		return nil
	})
	fmt.Println()
	return nil
}

// WithTheme sets the theme on a text field.
func (t *Text) WithTheme(theme *Theme) Field {
	if t.theme != nil {
		return t
	}
	t.theme = theme
	return t
}

// WithKeyMap sets the keymap on a text field.
func (t *Text) WithKeyMap(k *KeyMap) Field {
	t.keymap = k.Text
	t.textarea.KeyMap.InsertNewline.SetKeys(t.keymap.NewLine.Keys()...)
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
	t.textarea.SetWidth(width - t.activeStyles().Base.GetHorizontalFrameSize())
	return t
}

// WithHeight sets the height of the text field.
func (t *Text) WithHeight(height int) Field {
	adjust := 0
	if t.title != "" {
		adjust++
	}
	if t.description != "" {
		adjust++
	}
	t.textarea.SetHeight(height - t.activeStyles().Base.GetVerticalFrameSize() - adjust)
	return t
}

// WithPosition sets the position information of the text field.
func (t *Text) WithPosition(p FieldPosition) Field {
	t.keymap.Prev.SetEnabled(!p.IsFirst())
	t.keymap.Next.SetEnabled(!p.IsLast())
	t.keymap.Submit.SetEnabled(p.IsLast())
	return t
}

// GetKey returns the key of the field.
func (t *Text) GetKey() string {
	return t.key
}

// GetValue returns the value of the field.
func (t *Text) GetValue() any {
	return *t.value
}
