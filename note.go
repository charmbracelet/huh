package huh

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

// Note is a form note field.
type Note struct {
	body string

	style        *NoteStyle
	blurredStyle NoteStyle
	focusedStyle NoteStyle
}

// NewNote creates a new note field.
func NewNote() *Note {
	focused, blurred := DefaultNoteStyles()
	return &Note{
		focusedStyle: focused,
		blurredStyle: blurred,
	}
}

// Body sets the title of the select field.
func (n *Note) Body(body string) *Note {
	n.body = body
	return n
}

// Styles sets the styles of the select field.
func (n *Note) Styles(focused, blurred NoteStyle) *Note {
	n.blurredStyle = blurred
	n.focusedStyle = focused
	return n
}

// Focus focuses the select field.
func (n *Note) Focus() tea.Cmd {
	n.style = &n.focusedStyle
	return nil
}

// Blur blurs the select field.
func (n *Note) Blur() tea.Cmd {
	n.style = &n.blurredStyle
	return nil
}

// Init initializes the select field.
func (n *Note) Init() tea.Cmd {
	n.style = &n.blurredStyle
	return nil
}

// Update updates the select field.
func (n *Note) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "enter":
			return n, nextField
		}
		return n, nil
	}
	return n, nil
}

// View renders the select field.
func (n *Note) View() string {
	md, _ := glamour.Render(n.body, "auto")
	return n.style.Base.Render(md + n.style.Next.Render("Next"))
}

// Run runs an accessible select field.
func (n *Note) Run() {
	md, _ := glamour.Render(n.body, "auto")
	fmt.Println(md)
}
