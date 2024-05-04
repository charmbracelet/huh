package huh

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"
)

// Input is a form input field.
type Input struct {
	value *string
	key   string

	// customization
	title       string
	description string
	inline      bool

	// error handling
	validate func(string) error
	err      error

	// model
	textinput textinput.Model

	// state
	focused bool

	// options
	width      int
	height     int
	accessible bool
	theme      *Theme
	keymap     InputKeyMap
}

// NewInput returns a new input field.
func NewInput() *Input {
	input := textinput.New()

	i := &Input{
		value:     new(string),
		textinput: input,
		validate:  func(string) error { return nil },
	}

	return i
}

// Value sets the value of the input field.
func (i *Input) Value(value *string) *Input {
	i.value = value
	i.textinput.SetValue(*value)
	return i
}

// Key sets the key of the input field.
func (i *Input) Key(key string) *Input {
	i.key = key
	return i
}

// Title sets the title of the input field.
func (i *Input) Title(title string) *Input {
	i.title = title
	return i
}

// Description sets the description of the input field.
func (i *Input) Description(description string) *Input {
	i.description = description
	return i
}

// Prompt sets the prompt of the input field.
func (i *Input) Prompt(prompt string) *Input {
	i.textinput.Prompt = prompt
	return i
}

// CharLimit sets the character limit of the input field.
func (i *Input) CharLimit(charlimit int) *Input {
	i.textinput.CharLimit = charlimit
	return i
}

// Suggestions sets the suggestions to display for autocomplete in the input
// field.
func (i *Input) Suggestions(suggestions []string) *Input {
	i.textinput.ShowSuggestions = len(suggestions) > 0
	i.textinput.KeyMap.AcceptSuggestion.SetEnabled(len(suggestions) > 0)
	i.textinput.SetSuggestions(suggestions)
	return i
}

// EchoMode sets the input behavior of the text Input field.
type EchoMode textinput.EchoMode

const (
	// EchoNormal displays text as is.
	// This is the default behavior.
	EchoModeNormal EchoMode = EchoMode(textinput.EchoNormal)

	// EchoPassword displays the EchoCharacter mask instead of actual characters.
	// This is commonly used for password fields.
	EchoModePassword EchoMode = EchoMode(textinput.EchoPassword)

	// EchoNone displays nothing as characters are entered.
	// This is commonly seen for password fields on the command line.
	EchoModeNone EchoMode = EchoMode(textinput.EchoNone)
)

// EchoMode sets the echo mode of the input.
func (i *Input) EchoMode(mode EchoMode) *Input {
	i.textinput.EchoMode = textinput.EchoMode(mode)
	return i
}

// Password sets whether or not to hide the input while the user is typing.
//
// Deprecated: use EchoMode(EchoPassword) instead.
func (i *Input) Password(password bool) *Input {
	if password {
		i.textinput.EchoMode = textinput.EchoPassword
	} else {
		i.textinput.EchoMode = textinput.EchoNormal
	}
	return i
}

// Placeholder sets the placeholder of the text input.
func (i *Input) Placeholder(str string) *Input {
	i.textinput.Placeholder = str
	return i
}

// Inline sets whether the title and input should be on the same line.
func (i *Input) Inline(inline bool) *Input {
	i.inline = inline
	return i
}

// Validate sets the validation function of the input field.
func (i *Input) Validate(validate func(string) error) *Input {
	i.validate = validate
	return i
}

// Error returns the error of the input field.
func (i *Input) Error() error {
	return i.err
}

// Skip returns whether the input should be skipped or should be blocking.
func (*Input) Skip() bool {
	return false
}

// Zoom returns whether the input should be zoomed.
func (*Input) Zoom() bool {
	return false
}

// Focus focuses the input field.
func (i *Input) Focus() tea.Cmd {
	i.focused = true
	return i.textinput.Focus()
}

// Blur blurs the input field.
func (i *Input) Blur() tea.Cmd {
	i.focused = false
	*i.value = i.textinput.Value()
	i.textinput.Blur()
	i.err = i.validate(*i.value)
	return nil
}

// KeyBinds returns the help message for the input field.
func (i *Input) KeyBinds() []key.Binding {
	if i.textinput.ShowSuggestions {
		return []key.Binding{i.keymap.AcceptSuggestion, i.keymap.Prev, i.keymap.Submit, i.keymap.Next}
	}
	return []key.Binding{i.keymap.Prev, i.keymap.Submit, i.keymap.Next}
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
		i.err = nil

		switch {
		case key.Matches(msg, i.keymap.Prev):
			value := i.textinput.Value()
			i.err = i.validate(value)
			if i.err != nil {
				return i, nil
			}
			cmds = append(cmds, PrevField)
		case key.Matches(msg, i.keymap.Next, i.keymap.Submit):
			value := i.textinput.Value()
			i.err = i.validate(value)
			if i.err != nil {
				return i, nil
			}
			cmds = append(cmds, NextField)
		}
	}

	return i, tea.Batch(cmds...)
}

func (i *Input) activeStyles() *FieldStyles {
	theme := i.theme
	if theme == nil {
		theme = ThemeCharm()
	}
	if i.focused {
		return &theme.Focused
	}
	return &theme.Blurred
}

// View renders the input field.
func (i *Input) View() string {
	styles := i.activeStyles()

	// NB: since the method is on a pointer receiver these are being mutated.
	// Because this runs on every render this shouldn't matter in practice,
	// however.
	i.textinput.PlaceholderStyle = styles.TextInput.Placeholder
	i.textinput.PromptStyle = styles.TextInput.Prompt
	i.textinput.Cursor.Style = styles.TextInput.Cursor
	i.textinput.TextStyle = styles.TextInput.Text

	var sb strings.Builder
	if i.title != "" {
		sb.WriteString(styles.Title.Render(i.title))
		if !i.inline {
			sb.WriteString("\n")
		}
	}
	if i.description != "" {
		sb.WriteString(styles.Description.Render(i.description))
		if !i.inline {
			sb.WriteString("\n")
		}
	}
	sb.WriteString(i.textinput.View())

	return styles.Base.Render(sb.String())
}

// Run runs the input field in accessible mode.
func (i *Input) Run() error {
	if i.accessible {
		return i.runAccessible()
	}
	return i.run()
}

// run runs the input field.
func (i *Input) run() error {
	return Run(i)
}

// runAccessible runs the input field in accessible mode.
func (i *Input) runAccessible() error {
	styles := i.activeStyles()
	fmt.Println(styles.Title.Render(i.title))
	fmt.Println()
	*i.value = accessibility.PromptString("Input: ", i.validate)
	fmt.Println(styles.SelectedOption.Render("Input: " + *i.value + "\n"))
	return nil
}

// WithKeyMap sets the keymap on an input field.
func (i *Input) WithKeyMap(k *KeyMap) Field {
	i.keymap = k.Input
	i.textinput.KeyMap.AcceptSuggestion = i.keymap.AcceptSuggestion
	return i
}

// WithAccessible sets the accessible mode of the input field.
func (i *Input) WithAccessible(accessible bool) Field {
	i.accessible = accessible
	return i
}

// WithTheme sets the theme of the input field.
func (i *Input) WithTheme(theme *Theme) Field {
	if i.theme != nil {
		return i
	}
	i.theme = theme
	return i
}

// WithWidth sets the width of the input field.
func (i *Input) WithWidth(width int) Field {
	styles := i.activeStyles()
	i.width = width
	frameSize := styles.Base.GetHorizontalFrameSize()
	promptWidth := lipgloss.Width(i.textinput.PromptStyle.Render(i.textinput.Prompt))
	titleWidth := lipgloss.Width(styles.Title.Render(i.title))
	descriptionWidth := lipgloss.Width(styles.Description.Render(i.description))
	i.textinput.Width = width - frameSize - promptWidth - 1
	if i.inline {
		i.textinput.Width -= titleWidth
		i.textinput.Width -= descriptionWidth
	}
	return i
}

// WithHeight sets the height of the input field.
func (i *Input) WithHeight(height int) Field {
	i.height = height
	return i
}

// WithPosition sets the position of the input field.
func (i *Input) WithPosition(p FieldPosition) Field {
	i.keymap.Prev.SetEnabled(!p.IsFirst())
	i.keymap.Next.SetEnabled(!p.IsLast())
	i.keymap.Submit.SetEnabled(p.IsLast())
	return i
}

// GetKey returns the key of the field.
func (i *Input) GetKey() string {
	return i.key
}

// GetValue returns the value of the field.
func (i *Input) GetValue() any {
	return *i.value
}
