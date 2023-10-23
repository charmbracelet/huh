package huh

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"
)

// TextareaStyle is the style of the textarea field.
type TextStyle struct {
	Base        lipgloss.Style
	Title       lipgloss.Style
	Error       lipgloss.Style
	Description lipgloss.Style
	Help        lipgloss.Style
	textarea.Style
}

// DefaultTextStyles returns the default focused style of the text field.
func DefaultTextStyles() (TextStyle, TextStyle) {
	f, b := textarea.DefaultStyles()

	f.CursorLine = lipgloss.NewStyle()
	b.CursorLine = lipgloss.NewStyle()

	focused := TextStyle{
		Base:        lipgloss.NewStyle().Border(lipgloss.ThickBorder(), false).BorderLeft(true).PaddingLeft(1).MarginBottom(1).BorderForeground(lipgloss.Color("8")),
		Title:       lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		Description: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Error:       lipgloss.NewStyle().Foreground(lipgloss.Color("9")),
		Help:        lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Style:       f,
	}
	blurred := TextStyle{
		Base:        lipgloss.NewStyle().Border(lipgloss.HiddenBorder(), false).BorderLeft(true).PaddingLeft(1).MarginBottom(1),
		Title:       lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Error:       lipgloss.NewStyle().Foreground(lipgloss.Color("9")),
		Description: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Help:        lipgloss.NewStyle().Foreground(lipgloss.Color("0")),
		Style:       b,
	}

	return focused, blurred
}

// Text is a form text field.
type Text struct {
	value        *string
	title        string
	required     bool
	textarea     textarea.Model
	validate     func(string) error
	err          error
	style        *TextStyle
	focusedStyle TextStyle
	blurredStyle TextStyle
}

// NewText returns a new text field.
func NewText() *Text {
	text := textarea.New()
	text.ShowLineNumbers = false
	text.Prompt = ""

	f, b := DefaultTextStyles()

	t := &Text{
		value:        new(string),
		textarea:     text,
		style:        &b,
		focusedStyle: f,
		blurredStyle: b,
		validate:     func(s string) error { return nil },
	}

	t.updateTextareaStyle()

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

// Validate sets the validation function of the text field.
func (t *Text) Validate(validate func(string) error) *Text {
	t.validate = validate
	return t
}

// Error returns the error of the text field.
func (t *Text) Error() error {
	return t.err
}

// updateTextareaStyle updates the style of the textarea.
func (t *Text) updateTextareaStyle() {
	t.textarea.FocusedStyle = t.focusedStyle.Style
	t.textarea.BlurredStyle = t.blurredStyle.Style
}

// Focus focuses the text field.
func (t *Text) Focus() tea.Cmd {
	t.style = &t.focusedStyle
	cmd := t.textarea.Focus()
	return cmd
}

// Blur blurs the text field.
func (t *Text) Blur() tea.Cmd {
	*t.value = t.textarea.Value()
	t.style = &t.blurredStyle
	t.textarea.Blur()
	t.err = t.validate(t.textarea.Value())
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
		default:
			t.err = nil
		}
	}

	return t, tea.Batch(cmds...)
}

// View renders the text field.
func (t *Text) View() string {
	var sb strings.Builder
	sb.WriteString(t.style.Title.Render(t.title))
	if t.err != nil {
		sb.WriteString(t.style.Error.Render(" * "))
	}
	sb.WriteString("\n")
	sb.WriteString(t.textarea.View())
	sb.WriteString("\n")
	sb.WriteString(t.style.Help.Render("tab • continue"))

	return t.style.Base.Render(sb.String())
}

func (t *Text) Run() {
	fmt.Println(t.style.Title.Render(t.title))
	*t.value = accessibility.PromptString("> ")
	fmt.Println()
}
