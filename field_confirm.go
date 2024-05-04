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
	value *bool
	key   string

	// customization
	title       string
	description string
	affirmative string
	negative    string

	// error handling
	validate func(bool) error
	err      error

	// state
	focused bool

	// options
	width      int
	height     int
	inline     bool
	accessible bool
	theme      *Theme
	keymap     ConfirmKeyMap
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

// Skip returns whether the confirm should be skipped or should be blocking.
func (*Confirm) Skip() bool {
	return false
}

// Zoom returns whether the input should be zoomed.
func (*Confirm) Zoom() bool {
	return false
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

// Key sets the key of the confirm field.
func (c *Confirm) Key(key string) *Confirm {
	c.key = key
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

// Inline sets whether the field should be inline.
func (c *Confirm) Inline(inline bool) *Confirm {
	c.inline = inline
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

// KeyBinds returns the help message for the confirm field.
func (c *Confirm) KeyBinds() []key.Binding {
	return []key.Binding{c.keymap.Toggle, c.keymap.Prev, c.keymap.Submit, c.keymap.Next}
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
			if c.negative == "" {
				break
			}
			v := !*c.value
			*c.value = v
		case key.Matches(msg, c.keymap.Prev):
			cmds = append(cmds, PrevField)
		case key.Matches(msg, c.keymap.Next, c.keymap.Submit):
			cmds = append(cmds, NextField)
		}
	}

	return c, tea.Batch(cmds...)
}

func (c *Confirm) activeStyles() *FieldStyles {
	theme := c.theme
	if theme == nil {
		theme = ThemeCharm()
	}
	if c.focused {
		return &theme.Focused
	}
	return &theme.Blurred
}

// View renders the confirm field.
func (c *Confirm) View() string {
	styles := c.activeStyles()

	var sb strings.Builder
	sb.WriteString(styles.Title.Render(c.title))
	if c.err != nil {
		sb.WriteString(styles.ErrorIndicator.String())
	}

	description := styles.Description.Render(c.description)

	if !c.inline && c.description != "" {
		sb.WriteString("\n")
	}
	sb.WriteString(description)

	if !c.inline {
		sb.WriteString("\n")
		sb.WriteString("\n")
	}

	var negative string
	var affirmative string
	if c.negative != "" {
		if *c.value {
			affirmative = styles.FocusedButton.Render(c.affirmative)
			negative = styles.BlurredButton.Render(c.negative)
		} else {
			affirmative = styles.BlurredButton.Render(c.affirmative)
			negative = styles.FocusedButton.Render(c.negative)
		}
	} else {
		affirmative = styles.FocusedButton.Render(c.affirmative)
	}

	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Center, affirmative, negative))
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
	styles := c.activeStyles()
	fmt.Println(styles.Title.Render(c.title))
	fmt.Println()
	*c.value = accessibility.PromptBool()
	fmt.Println(styles.SelectedOption.Render("Chose: "+c.String()) + "\n")
	return nil
}

func (c *Confirm) String() string {
	if *c.value {
		return c.affirmative
	}
	return c.negative
}

// WithTheme sets the theme of the confirm field.
func (c *Confirm) WithTheme(theme *Theme) Field {
	if c.theme != nil {
		return c
	}
	c.theme = theme
	return c
}

// WithKeyMap sets the keymap of the confirm field.
func (c *Confirm) WithKeyMap(k *KeyMap) Field {
	c.keymap = k.Confirm
	return c
}

// WithAccessible sets the accessible mode of the confirm field.
func (c *Confirm) WithAccessible(accessible bool) Field {
	c.accessible = accessible
	return c
}

// WithWidth sets the width of the confirm field.
func (c *Confirm) WithWidth(width int) Field {
	c.width = width
	return c
}

// WithHeight sets the height of the confirm field.
func (c *Confirm) WithHeight(height int) Field {
	c.height = height
	return c
}

// WithPosition sets the position of the confirm field.
func (c *Confirm) WithPosition(p FieldPosition) Field {
	c.keymap.Prev.SetEnabled(!p.IsFirst())
	c.keymap.Next.SetEnabled(!p.IsLast())
	c.keymap.Submit.SetEnabled(p.IsLast())
	return c
}

// GetKey returns the key of the field.
func (c *Confirm) GetKey() string {
	return c.key
}

// GetValue returns the value of the field.
func (c *Confirm) GetValue() any {
	return *c.value
}
