package huh

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

// Note is a form note field.
type Note struct {
	body string

	showNextButton bool
	style          *NoteStyle
	blurredStyle   NoteStyle
	focusedStyle   NoteStyle
}

// NewNote creates a new note field.
func NewNote() *Note {
	focused, blurred := DefaultNoteStyles()
	return &Note{
		showNextButton: false,
		style:          &blurred,
		focusedStyle:   focused,
		blurredStyle:   blurred,
	}
}

// Body sets the title of the select field.
func (n *Note) Body(body string) *Note {
	n.body = body
	return n
}

// Next sets whether to show the next button.
func (n *Note) Next(show bool) *Note {
	n.showNextButton = show
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
	return nil
}

// Update updates the select field.
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

// View renders the select field.
func (n *Note) View() string {
	var sb strings.Builder
	md, _ := glamour.Render(n.body, "auto")
	sb.WriteString(md)
	if n.showNextButton {
		sb.WriteString(n.style.Next.Render("Next"))
	}
	return n.style.Base.Render(sb.String())
}

// Run runs an accessible select field.
func (n *Note) Run() {
	md, _ := glamour.Render(n.body, "auto")
	fmt.Println(md)
}
