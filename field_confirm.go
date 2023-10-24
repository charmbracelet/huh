package huh

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"
)

// Confirm is a form confirm field.
type Confirm struct {
	value       *bool
	title       string
	description string

	validate func(bool) error
	err      error

	affirmative string
	negative    string

	focused    bool
	accessible bool
	theme      *Theme
	keymap     *ConfirmKeyMap
}

// NewConfirm returns a new confirm field.
func NewConfirm() *Confirm {
	return &Confirm{
		value:       new(bool),
		affirmative: "Yes",
		negative:    "No",
		validate:    func(bool) error { return nil },
	}
}

// Validate sets the validation function of the confirm field.
func (c *Confirm) Validate(validate func(bool) error) *Confirm {
	c.validate = validate
	return c
}

// Error returns the error of the confirm field.
func (c *Confirm) Error() error {
	return c.err
}

// Affirmative sets the affirmative value of the confirm field.
func (c *Confirm) Affirmative(affirmative string) *Confirm {
	c.affirmative = affirmative
	return c
}

// Negative sets the negative value of the confirm field.
func (c *Confirm) Negative(negative string) *Confirm {
	c.negative = negative
	return c
}

// Value sets the value of the confirm field.
func (c *Confirm) Value(value *bool) *Confirm {
	c.value = value
	return c
}

// Title sets the title of the confirm field.
func (c *Confirm) Title(title string) *Confirm {
	c.title = title
	return c
}

// Description sets the description of the confirm field.
func (c *Confirm) Description(description string) *Confirm {
	c.description = description
	return c
}

// Focus focuses the confirm field.
func (c *Confirm) Focus() tea.Cmd {
	c.focused = true
	return nil
}

// Blur blurs the confirm field.
func (c *Confirm) Blur() tea.Cmd {
	c.focused = false
	c.err = c.validate(*c.value)
	return nil
}

// KeyMap sets the keymap of the confirm field.
func (c *Confirm) KeyMap(k *KeyMap) Field {
	c.keymap = &k.Confirm
	return c
}

// KeyBinds returns the help message for the confirm field.
func (c *Confirm) KeyBinds() []key.Binding {
	return []key.Binding{c.keymap.Toggle, c.keymap.Next, c.keymap.Prev}
}

// Accessible sets the accessible mode of the confirm field.
func (c *Confirm) Accessible(accessible bool) Field {
	c.accessible = accessible
	return c
}

// Init initializes the confirm field.
func (c *Confirm) Init() tea.Cmd {
	return nil
}

// Update updates the confirm field.
func (c *Confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:

		c.err = nil

		switch {
		case key.Matches(msg, c.keymap.Toggle):
			*c.value = !*c.value
		case key.Matches(msg, c.keymap.Prev):
			cmds = append(cmds, prevField)
		case key.Matches(msg, c.keymap.Next):
			cmds = append(cmds, nextField)
		}
	}

	return c, tea.Batch(cmds...)
}

// View renders the confirm field.
func (c *Confirm) View() string {
	styles := c.theme.Blurred
	if c.focused {
		styles = c.theme.Focused
	}

	var sb strings.Builder
	sb.WriteString(styles.Title.Render(c.title))
	if c.err != nil {
		sb.WriteString(styles.ErrorIndicator.String())
	}
	if c.description != "" {
		sb.WriteString("\n")
		sb.WriteString(styles.Description.Render(c.description))
	}
	sb.WriteString("\n")
	sb.WriteString("\n")

	if *c.value {
		sb.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Center,
			styles.FocusedButton.Render(c.affirmative),
			styles.BlurredButton.Render(c.negative),
		))
	} else {
		sb.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Center,
			styles.BlurredButton.Render(c.affirmative),
			styles.FocusedButton.Render(c.negative),
		))
	}
	return styles.Base.Render(sb.String())
}

// Run runs the confirm field in accessible mode.
func (c *Confirm) Run() error {
	if c.accessible {
		return c.runAccessible()
	}
	return Run(c)
}

// runAccessible runs the confirm field in accessible mode.
func (c *Confirm) runAccessible() error {
	fmt.Println(c.theme.Focused.Title.Render(c.title))
	choice := accessibility.PromptBool()
	*c.value = choice
	if choice {
		fmt.Println("Selected: " + c.affirmative)
	} else {
		fmt.Println("Selected: " + c.negative)
	}
	fmt.Println()
	return nil
}

func (c *Confirm) Theme(theme *Theme) Field {
	c.theme = theme
	return c
}
