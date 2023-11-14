package huh

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

func init() {
	// XXX: For now, unset padding on default glamour styles.
	glamour.DarkStyleConfig.Document.BlockPrefix = ""
	glamour.DarkStyleConfig.Document.Margin = pointerTo[uint](0)
	glamour.LightStyleConfig.Document.BlockPrefix = ""
	glamour.LightStyleConfig.Document.Margin = pointerTo[uint](0)
}

// Note is a form note field.
type Note struct {
	// customization
	title       string
	description string

	// state
	showNextButton bool
	focused        bool
	renderer       *glamour.TermRenderer

	// options
	width      int
	accessible bool
	theme      *Theme
	keymap     *NoteKeyMap
}

// NewNote creates a new note field.
func NewNote() *Note {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithEmoji(),
		glamour.WithWordWrap(0),
	)

	if err != nil {
		r, _ = glamour.NewTermRenderer()
	}

	return &Note{
		showNextButton: false,
		renderer:       r,
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

// KeyBinds returns the help message for the note field.
func (n *Note) KeyBinds() []key.Binding {
	return []key.Binding{n.keymap.Next}
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
		case key.Matches(msg, n.keymap.Next):
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
		sb   strings.Builder
		body string
	)

	if n.title != "" {
		body = fmt.Sprintf("# %s\n", n.title)
	}

	body += n.description

	md, _ := n.renderer.Render(body)
	sb.WriteString(md)
	if n.showNextButton {
		sb.WriteString(styles.Next.Render("Next"))
	}
	return styles.Base.Render(sb.String())
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
		body = fmt.Sprintf("# %s\n", n.title)
	}

	body += n.description

	md, _ := n.renderer.Render(body)
	fmt.Println(n.theme.Blurred.Base.Render(strings.TrimSpace(md)))
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
	n.keymap = &k.Note
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

// GetValue satisfies the Field interface, notes do not have values.
func (n *Note) GetValue() any {
	return nil
}

// GetKey satisfies the Field interface, notes do not have keys.
func (n *Note) GetKey() string {
	return ""
}

// pointerTo returns a pointer to a value.
func pointerTo[T any](v T) *T {
	return &v
}
