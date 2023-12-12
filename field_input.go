package huh

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
)

// Input is a form input field.
type Input struct {
	value *string
	key   string

	// customization
	title       string
	description string
	inline      bool
	charlimit   int

	// error handling
	validate func(string) error
	err      error

	// model
	textinput textinput.Model

	// state
	focused bool

	// options
	width      int
	accessible bool
	theme      *Theme
	keymap     *InputKeyMap
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
	i.charlimit = charlimit
	return i
}

// Password sets whether or not to hide the input while the user is typing.
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
	return []key.Binding{i.keymap.Next, i.keymap.Prev}
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
			cmds = append(cmds, prevField)
		case key.Matches(msg, i.keymap.Next):
			value := i.textinput.Value()
			i.err = i.validate(value)
			if i.err != nil {
				return i, nil
			}
			cmds = append(cmds, nextField)
		}
	}

	return i, tea.Batch(cmds...)
}

// View renders the input field.
func (i *Input) View() string {
	styles := i.theme.Blurred
	if i.focused {
		styles = i.theme.Focused
	}

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
	fmt.Println(i.theme.Blurred.Base.Render(i.theme.Focused.Title.Render(i.title)))
	fmt.Println()
	*i.value = accessibility.PromptString("Input: ", i.validate)
	fmt.Println(i.theme.Focused.SelectedOption.Render("Input: " + *i.value + "\n"))
	return nil
}

// WithKeyMap sets the keymap on an input field.
func (i *Input) WithKeyMap(k *KeyMap) Field {
	i.keymap = &k.Input
	return i
}

// WithAccessible sets the accessible mode of the input field.
func (i *Input) WithAccessible(accessible bool) Field {
	i.accessible = accessible
	return i
}

// WithTheme sets the theme of the input field.
func (i *Input) WithTheme(theme *Theme) Field {
	i.theme = theme
	return i
}

// WithWidth sets the width of the input field.
func (i *Input) WithWidth(width int) Field {
	i.width = width
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
