package huh

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Note is a form note field.
type Note struct {
	// customization
	title       string
	description string

	// state
	showNextButton bool
	focused        bool

	// options
	skip       bool
	width      int
	height     int
	accessible bool
	theme      *Theme
	keymap     NoteKeyMap
}

// NewNote creates a new note field.
func NewNote() *Note {
	return &Note{
		showNextButton: false,
		theme:          ThemeCharm(),
		skip:           true,
	}
}

// Title sets the title of the note field.
func (n *Note) Title(title string) *Note {
	n.title = title
	return n
}

// Description sets the description of the note field.
func (n *Note) Description(description string) *Note {
	n.description = description
	return n
}

// Next sets whether to show the next button.
func (n *Note) Next(show bool) *Note {
	n.showNextButton = show
	return n
}

// Focus focuses the note field.
func (n *Note) Focus() tea.Cmd {
	n.focused = true
	return nil
}

// Blur blurs the note field.
func (n *Note) Blur() tea.Cmd {
	n.focused = false
	return nil
}

// Error returns the error of the note field.
func (n *Note) Error() error {
	return nil
}

// Skip returns whether the note should be skipped or should be blocking.
func (n *Note) Skip() bool {
	return n.skip
}

// KeyBinds returns the help message for the note field.
func (n *Note) KeyBinds() []key.Binding {
	return []key.Binding{n.keymap.Prev, n.keymap.Submit, n.keymap.Next}
}

// Init initializes the note field.
func (n *Note) Init() tea.Cmd {
	return nil
}

// Update updates the note field.
func (n *Note) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, n.keymap.Prev):
			return n, prevField
		case key.Matches(msg, n.keymap.Next, n.keymap.Submit):
			return n, nextField
		}
		return n, nextField
	}
	return n, nil
}

// View renders the note field.
func (n *Note) View() string {
	styles := n.theme.Blurred
	if n.focused {
		styles = n.theme.Focused
	}

	var (
		sb    strings.Builder
		title string
	)

	if n.title != "" {
		title = n.title
	}
	sb.WriteString(styles.NoteTitle.Render(title))
	if n.description != "" {
		sb.WriteString("\n")
		sb.WriteString(render(n.description))
	}
	if n.showNextButton {
		sb.WriteString(styles.Next.Render("Next"))
	}
	return styles.Card.Render(sb.String())
}

// Run runs the note field.
func (n *Note) Run() error {
	if n.accessible {
		return n.runAccessible()
	}
	return Run(n)
}

// runAccessible runs an accessible note field.
func (n *Note) runAccessible() error {
	var body string

	if n.title != "" {
		body = n.title + "\n\n"
	}

	body += n.description

	fmt.Println(n.theme.Blurred.Base.Render(body))
	fmt.Println()
	return nil
}

// WithTheme sets the theme on a note field.
func (n *Note) WithTheme(theme *Theme) Field {
	n.theme = theme
	return n
}

// WithKeyMap sets the keymap on a note field.
func (n *Note) WithKeyMap(k *KeyMap) Field {
	n.keymap = k.Note
	return n
}

// WithAccessible sets the accessible mode of the note field.
func (n *Note) WithAccessible(accessible bool) Field {
	n.accessible = accessible
	return n
}

// WithWidth sets the width of the note field.
func (n *Note) WithWidth(width int) Field {
	n.width = width
	return n
}

// WithHeight sets the height of the note field.
func (n *Note) WithHeight(height int) Field {
	n.height = height
	return n
}

// WithPosition sets the position information of the note field.
func (n *Note) WithPosition(p FieldPosition) Field {
	// if the note is the only field on the screen,
	// we shouldn't skip the entire group.
	if p.Field == p.FirstField && p.Field == p.LastField {
		n.skip = false
	}
	n.keymap.Prev.SetEnabled(!p.IsFirst())
	n.keymap.Next.SetEnabled(!p.IsLast())
	n.keymap.Submit.SetEnabled(p.IsLast())
	return n
}

// GetValue satisfies the Field interface, notes do not have values.
func (n *Note) GetValue() any {
	return nil
}

// GetKey satisfies the Field interface, notes do not have keys.
func (n *Note) GetKey() string {
	return ""
}

func render(input string) string {
	var result strings.Builder
	var italic, bold, codeblock bool

	for _, char := range input {
		switch char {
		case '_':
			if !italic {
				result.WriteString("\033[3m")
				italic = true
			} else {
				result.WriteString("\033[0m")
				italic = false
			}
		case '*':
			if !bold {
				result.WriteString("\033[1m")
				bold = true
			} else {
				result.WriteString("\033[0m")
				bold = false
			}
		case '`':
			if !codeblock {
				result.WriteString("\033[0;37;40m")
				result.WriteString(" ")
				codeblock = true
			} else {
				result.WriteString(" ")
				result.WriteString("\033[0m")
				codeblock = false
			}
		default:
			result.WriteRune(char)
		}
	}

	// Reset any open formatting
	result.WriteString("\033[0m")

	return result.String()
}
