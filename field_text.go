package huh

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
)

// Text is a form text field.
type Text struct {
	value    *string
	title    string
	required bool
	textarea textarea.Model
	focused  bool
	theme    *Theme
}

// NewText returns a new text field.
func NewText() *Text {
	text := textarea.New()
	text.ShowLineNumbers = false
	text.Prompt = ""

	t := &Text{
		value:    new(string),
		textarea: text,
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

// Required sets the text field as required.
func (t *Text) Required(required bool) *Text {
	t.required = required
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
	return nil
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
		switch msg.String() {
		case "tab", "ctrl+d":
			cmds = append(cmds, nextField)
		case "shift+tab":
			cmds = append(cmds, prevField)
		}
	}

	return t, tea.Batch(cmds...)
}

// View renders the text field.
func (t *Text) View() string {
	styles := t.theme.Blurred
	if t.focused {
		styles = t.theme.Focused
	}

	var sb strings.Builder
	sb.WriteString(styles.Title.Render(t.title))
	sb.WriteString("\n")
	sb.WriteString(t.textarea.View())
	sb.WriteString("\n")
	sb.WriteString(styles.Help.Render("tab • continue"))

	return styles.Base.Render(sb.String())
}

func (t *Text) Run() {
	fmt.Println(t.theme.Focused.Title.Render(t.title))
	*t.value = accessibility.PromptString("> ")
	fmt.Println()
}

func (t *Text) setTheme(theme *Theme) {
	t.theme = theme
}
