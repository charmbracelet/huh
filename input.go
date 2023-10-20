package huh

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
)

// Input is a form input field.
type Input struct {
	value        *string
	title        string
	required     bool
	charlimit    int
	textinput    textinput.Model
	style        *InputStyle
	focusedStyle InputStyle
	blurredStyle InputStyle
}

// NewInput returns a new input field.
func NewInput() *Input {
	input := textinput.New()

	f, b := DefaultInputStyles()

	i := &Input{
		value:        new(string),
		textinput:    input,
		style:        &b,
		focusedStyle: f,
		blurredStyle: b,
	}

	i.updateTextinputStyle()

	return i
}

// Value sets the value of the input field.
func (i *Input) Value(value *string) *Input {
	i.value = value
	return i
}

// Title sets the title of the input field.
func (i *Input) Title(title string) *Input {
	i.title = title
	return i
}

// Required sets the input field as required.
func (i *Input) Required(required bool) *Input {
	i.required = required
	return i
}

// CharLimit sets the character limit of the input field.
func (i *Input) CharLimit(charlimit int) *Input {
	i.charlimit = charlimit
	return i
}

// Placeholder sets the placeholder of the text input.
func (i *Input) Placeholder(str string) *Input {
	i.textinput.Placeholder = str
	return i
}

func (i *Input) updateTextinputStyle() {
	i.textinput.PromptStyle = i.style.Prompt
	i.textinput.PlaceholderStyle = i.style.Placeholder
	i.textinput.TextStyle = i.style.Text
}

// Focus focuses the input field.
func (i *Input) Focus() tea.Cmd {
	i.style = &i.focusedStyle
	cmd := i.textinput.Focus()
	i.updateTextinputStyle()
	return cmd
}

// Blur blurs the input field.
func (i *Input) Blur() tea.Cmd {
	i.style = &i.blurredStyle
	i.textinput.Blur()
	i.updateTextinputStyle()
	return nil
}

// Init initializes the input field.
func (i *Input) Init() tea.Cmd {
	i.textinput.Blur()
	return nil
}

// Update updates the input field.
func (i *Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	i.textinput, cmd = i.textinput.Update(msg)
	cmds = append(cmds, cmd)
	*i.value = i.textinput.Value()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "shift+tab":
			cmds = append(cmds, prevField)
		case "enter", "tab":
			cmds = append(cmds, nextField)
		}
	}

	return i, tea.Batch(cmds...)
}

// View renders the input field.
func (i *Input) View() string {
	var sb strings.Builder

	sb.WriteString(i.style.Title.Render(i.title))
	sb.WriteString("\n")
	sb.WriteString(i.textinput.View())

	return i.style.Base.Render(sb.String())
}

// Run runs the input field in accessible mode.
func (i *Input) Run() {
	fmt.Println(i.style.Title.Render(i.title))
	*i.value = accessibility.PromptString("> ")
	fmt.Println()
}
