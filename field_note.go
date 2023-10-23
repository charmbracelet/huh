package huh

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// NoteStyle is the style of the Note field.
type NoteStyle struct {
	Base        lipgloss.Style
	Title       lipgloss.Style
	Description lipgloss.Style
	Next        lipgloss.Style
	Body        lipgloss.Style
}

// DefaultNoteStyles returns the default focused style of the Note field.
func DefaultNoteStyles() (NoteStyle, NoteStyle) {
	focused := NoteStyle{
		Base:        lipgloss.NewStyle().Border(lipgloss.ThickBorder(), false).Margin(1, 0).BorderForeground(lipgloss.Color("8")),
		Title:       lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Margin(1),
		Description: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Body:        lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Margin(1),
		Next:        lipgloss.NewStyle().Background(lipgloss.Color("5")).Foreground(lipgloss.Color("3")).Margin(1, 2).Padding(0, 1).Bold(true),
	}
	blurred := NoteStyle{
		Base:        lipgloss.NewStyle().Border(lipgloss.HiddenBorder(), false).BorderLeft(true).MarginBottom(1),
		Title:       lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Margin(1),
		Description: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Body:        lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Margin(1),
		Next:        lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Margin(0, 2).MarginBottom(1),
	}
	return focused, blurred
}

// Note is a form note field.
type Note struct {
	title       string
	description string

	showNextButton bool
	style          *NoteStyle
	blurredStyle   NoteStyle
	focusedStyle   NoteStyle
	theme          *Theme
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

// Styles sets the styles of the note field.
func (n *Note) Styles(focused, blurred NoteStyle) *Note {
	n.blurredStyle = blurred
	n.focusedStyle = focused
	return n
}

// Focus focuses the note field.
func (n *Note) Focus() tea.Cmd {
	n.style = &n.focusedStyle
	return nil
}

// Blur blurs the note field.
func (n *Note) Blur() tea.Cmd {
	n.style = &n.blurredStyle
	return nil
}

// Init initializes the note field.
func (n *Note) Init() tea.Cmd {
	return nil
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
	var sb strings.Builder

	var body string

	if n.title != "" {
		body = fmt.Sprintf("# %s\n", n.title)
	}

	body += n.description

	md, _ := glamour.Render(body, "auto")
	sb.WriteString(md)
	if n.showNextButton {
		sb.WriteString(n.style.Next.Render("Next"))
	}
	return n.style.Base.Render(strings.TrimSpace(sb.String()))
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

func (n *Note) setTheme(theme *Theme) {
	n.theme = theme
}
