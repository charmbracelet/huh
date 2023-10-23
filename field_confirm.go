package huh

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmStyle is the style of the confirm field.
type ConfirmStyle struct {
	Base        lipgloss.Style
	Title       lipgloss.Style
	Description lipgloss.Style
	Selected    lipgloss.Style
	Unselected  lipgloss.Style
}

// DefaultConfirmStyles returns the default focused style of the confirm field.
func DefaultConfirmStyles() (ConfirmStyle, ConfirmStyle) {
	focused := ConfirmStyle{
		Base:        lipgloss.NewStyle().Border(lipgloss.ThickBorder(), false).BorderLeft(true).PaddingLeft(1).MarginBottom(1).BorderForeground(lipgloss.Color("8")),
		Title:       lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		Description: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Selected:    lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Background(lipgloss.Color("4")).Padding(0, 2).MarginLeft(2),
		Unselected:  lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Background(lipgloss.Color("0")).Padding(0, 2).MarginLeft(2),
	}
	blurred := ConfirmStyle{
		Base:        lipgloss.NewStyle().Border(lipgloss.HiddenBorder(), false).BorderLeft(true).PaddingLeft(1).MarginBottom(1),
		Title:       lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Description: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Selected:    lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Background(lipgloss.Color("0")).Padding(0, 2).MarginLeft(2),
		Unselected:  lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Background(lipgloss.Color("0")).Padding(0, 2).MarginLeft(2),
	}
	return focused, blurred
}

// Confirm is a form confirm field.
type Confirm struct {
	value       *bool
	title       string
	description string
	required    bool

	affirmative string
	negative    string

	style        *ConfirmStyle
	focusedStyle ConfirmStyle
	blurredStyle ConfirmStyle

	theme *Theme
}

// NewConfirm returns a new confirm field.
func NewConfirm() *Confirm {
	f, b := DefaultConfirmStyles()
	return &Confirm{
		value:        new(bool),
		focusedStyle: f,
		blurredStyle: b,
		affirmative:  "Yes",
		negative:     "No",
	}
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

// Required sets the confirm field as required.
func (c *Confirm) Required(required bool) *Confirm {
	c.required = required
	return c
}

// Focus focuses the confirm field.
func (c *Confirm) Focus() tea.Cmd {
	c.style = &c.focusedStyle
	return nil
}

// Blur blurs the confirm field.
func (c *Confirm) Blur() tea.Cmd {
	c.style = &c.blurredStyle
	return nil
}

// Init initializes the confirm field.
func (c *Confirm) Init() tea.Cmd {
	c.style = &c.blurredStyle
	return nil
}

// Update updates the confirm field.
func (c *Confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			*c.value = true
		case "n", "N":
			*c.value = false
		case "h", "l", "left", "right":
			*c.value = !*c.value
		case "shift+tab":
			cmds = append(cmds, prevField)
		case "enter", "tab":
			cmds = append(cmds, nextField)
		}
	}

	return c, tea.Batch(cmds...)
}

// View renders the confirm field.
func (c *Confirm) View() string {
	style := c.style
	if style == nil {
		style = &c.blurredStyle
	}

	var sb strings.Builder
	sb.WriteString(c.style.Title.Render(c.title))
	if c.description != "" {
		sb.WriteString("\n")
		sb.WriteString(c.style.Description.Render(c.description))
	}
	sb.WriteString("\n")
	sb.WriteString("\n")

	if *c.value {
		sb.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Center,
			style.Selected.Render(c.affirmative),
			style.Unselected.Render(c.negative),
		))
	} else {
		sb.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Center,
			style.Unselected.Render(c.affirmative),
			style.Selected.Render(c.negative),
		))
	}
	return c.style.Base.Render(sb.String())
}

// Run runs the confirm field in accessible mode.
func (c *Confirm) Run() {
	fmt.Println(c.style.Title.Render(c.title))
	choice := accessibility.PromptBool()
	*c.value = choice
	if choice {
		fmt.Println("Selected: " + c.affirmative)
	} else {
		fmt.Println("Selected: " + c.negative)
	}
	fmt.Println()
}

func (c *Confirm) setTheme(theme *Theme) {
	c.theme = theme
}
