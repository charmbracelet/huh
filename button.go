package huh

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
)

// Button is a pressable button with an optional action
type Button struct {
	title       string
	description string
	inline      bool
	fn          func() error
	lastErr     error

	focused    bool
	accessible bool
	skip       bool
	width      int
	height     int
	theme      *Theme
	keymap     MultiSelectKeyMap
}

var _ Field = (*Button)(nil)

// NewButton initializes a new button with no values set
func NewButton() *Button {
	return &Button{}
}

// Init initializes the button
func (b *Button) Init() tea.Cmd {
	return nil
}

// Update handles interactions on the button
func (b *Button) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keymap.Prev):
			return b, PrevField
		case key.Matches(msg, b.keymap.Next):
			return b, NextField
		case key.Matches(msg, b.keymap.Submit):
			return b, nil
		case key.Matches(msg, b.keymap.Toggle):
			b.lastErr = b.fn()
			return b, nil
		}
	}

	return b, nil
}

// View renders the button
func (b *Button) View() string {
	styles := b.activeStyles()

	var builder strings.Builder

	title := b.title
	if title == "" {
		title = "<no title>"
	}

	if b.focused {
		builder.WriteString(styles.FocusedButton.Render(title))
	} else {
		builder.WriteString(styles.BlurredButton.Render(title))
	}

	if b.description != "" {
		if b.inline {
			builder.WriteString(" - ")
		} else {
			builder.WriteString("\n")
		}

		builder.WriteString(styles.Description.Render(b.description))
	}

	return styles.Base.Render(builder.String())
}

func (b *Button) activeStyles() *FieldStyles {
	theme := b.theme

	if theme == nil {
		theme = ThemeCharm()
	}

	if b.focused {
		return &theme.Focused
	}

	return &theme.Blurred
}

// Blur vlurs the button.
func (b *Button) Blur() tea.Cmd {
	b.focused = false

	return nil
}

// Focus focuses the button.
func (b *Button) Focus() tea.Cmd {
	b.focused = true

	return nil
}

// Error returns an error if one was returned from the button's execution
func (b *Button) Error() error {
	return b.lastErr
}

// Run runs the field individually.
func (b *Button) Run() error {
	if b.accessible {
		return b.runAccessible()
	}

	return Run(b)
}

func (b *Button) runAccessible() error {
	styles := b.activeStyles()
	fmt.Println(styles.Title.Render(b.title))
	fmt.Println()

	if accessibility.PromptBool() {
		fmt.Println(styles.SelectedOption.Render("Button pressed"))

		b.lastErr = b.fn()
		if b.lastErr != nil {
			fmt.Println(styles.ErrorMessage.Render("Error: " + b.lastErr.Error()))
		}
	}

	return nil
}

// Skip returns whether this input should be skipped or not.
func (b *Button) Skip() bool {
	return b.skip
}

// Zoom returns whether this button should be zoomed or not.
func (b *Button) Zoom() bool {
	return false
}

// KeyBinds returns help keybindings.
func (b *Button) KeyBinds() []key.Binding {
	return []key.Binding{b.keymap.Prev, b.keymap.Toggle, b.keymap.Submit, b.keymap.Next}
}

// WithTheme sets the theme of the button.
func (b *Button) WithTheme(theme *Theme) Field {
	if b.theme != nil {
		return b
	}

	b.theme = theme

	return b
}

// WithAccessible sets whether the field should run in accessible mode.
func (b *Button) WithAccessible(accessible bool) Field {
	b.accessible = accessible
	return b
}

// WithKeyMap sets the keymap on a field.
func (b *Button) WithKeyMap(keymap *KeyMap) Field {
	b.keymap = keymap.MultiSelect

	return b
}

// WithWidth sets the width of a field.
func (b *Button) WithWidth(width int) Field {
	b.width = width

	return b
}

// WithHeight sets the height of a field.
func (b *Button) WithHeight(height int) Field {
	b.height = height

	return b
}

// WithPosition tells the field the index of the group and position it is in.
func (b *Button) WithPosition(p FieldPosition) Field {
	if p.Field == p.FirstField && p.Field == p.LastField {
		b.skip = false
	}

	b.keymap.Prev.SetEnabled(!p.IsFirst())
	b.keymap.Next.SetEnabled(!p.IsLast())
	b.keymap.Submit.SetEnabled(b.fn != nil)

	return b
}

// GetKey returns the field's key.
func (b *Button) GetKey() string {
	return ""
}

// GetValue returns the field's value.
func (b *Button) GetValue() any {
	return nil
}

// Title sets the title of the button.
func (b *Button) Title(title string) *Button {
	b.title = title
	return b
}

// Description sets the button's description.
func (b *Button) Description(description string) *Button {
	b.description = description
	return b
}

// Action defines the button that will be activated when the button is preseed.
func (b *Button) Action(action func() error) *Button {
	b.fn = action
	return b
}

// Inline declares whether the button's description will be displayed inline.
func (b *Button) Inline(inline bool) *Button {
	b.inline = inline
	return b
}
