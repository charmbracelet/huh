package huh

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

// Note is a form note field.
type Note struct {
	title       string
	description string

	showNextButton bool
	focused        bool
	theme          *Theme
	keymap         *NoteKeyMap
}

// NewNote creates a new note field.
func NewNote() *Note {
	return &Note{
		showNextButton: false,
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

// Init initializes the note field.
func (n *Note) Init() tea.Cmd {
	return nil
}

// KeyMap sets the keymap on a note field.
func (n *Note) KeyMap(k *KeyMap) Field {
	n.keymap = &k.Note
	return n
}

// KeyBinds returns the help message for the note field.
func (n *Note) KeyBinds() []key.Binding {
	return []key.Binding{n.keymap.Next}
}

// Update updates the note field.
func (n *Note) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "shift+tab":
			return n, prevField
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
		sb   strings.Builder
		body string
	)

	if n.title != "" {
		body = fmt.Sprintf("# %s\n", n.title)
	}

	body += n.description

	md, _ := glamour.Render(body, "auto")
	sb.WriteString(md)
	if n.showNextButton {
		sb.WriteString(styles.Next.Render("Next"))
		sb.WriteString("\n")
	}
	return styles.Base.Render(sb.String())
}

// Run runs an accessible note field.
func (n *Note) Run() {
	var body string

	if n.title != "" {
		body = fmt.Sprintf("# %s\n", n.title)
	}

	body += n.description

	md, _ := glamour.Render(body, "auto")
	fmt.Println(strings.TrimSpace(md))
}

func (n *Note) Theme(theme *Theme) Field {
	n.theme = theme
	return n
}
