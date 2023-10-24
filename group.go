package huh

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/ordered"
)

// Group is a collection of fields that are displayed together with a page of
// the form. While a group is displayed the form completer can switch between
// fields in the group.
//
// If any of the fields in a group have errors, the form will not be able to
// progress to the next group.
type Group struct {
	// collection of fields
	fields []Field

	// information
	title       string
	description string

	// navigation
	current int

	// help
	showHelp bool
	help     help.Model

	// form options
	theme  *Theme
	keymap *KeyMap
}

// NewGroup returns a new group with the given fields.
func NewGroup(fields ...Field) *Group {
	return &Group{
		fields:   fields,
		current:  0,
		help:     help.New(),
		showHelp: true,
	}
}

// Title sets the group's title.
func (g *Group) Title(title string) *Group {
	g.title = title
	return g
}

// Description sets the group's description.
func (g *Group) Description(description string) *Group {
	g.description = description
	return g
}

// ShowHelp sets whether or not the group's help should be shown.
func (g *Group) ShowHelp(showHelp bool) *Group {
	g.showHelp = showHelp
	return g
}

// Theme sets the theme on a group.
func (g *Group) Theme(t *Theme) *Group {
	g.theme = t
	return g
}

// KeyMap sets the keymap on a group.
func (g *Group) KeyMap(k *KeyMap) *Group {
	g.keymap = k
	return g
}

// Errors returns the groups' fields' errors.
func (g *Group) Errors() []error {
	var errs []error
	for _, field := range g.fields {
		if err := field.Error(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// nextFieldMsg is a message to move to the next field,
//
// each field controls when to send this message such that it is able to use
// different key bindings or events to trigger group progression.
type nextFieldMsg struct{}

// prevFieldMsg is a message to move to the previous field.
//
// each field controls when to send this message such that it is able to use
// different key bindings or events to trigger group progression.
type prevFieldMsg struct{}

// nextField is the command to move to the next field.
func nextField() tea.Msg {
	return nextFieldMsg{}
}

// prevField is the command to move to the previous field.
func prevField() tea.Msg {
	return prevFieldMsg{}
}

// Init initializes the group.
func (g *Group) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, field := range g.fields {
		cmds = append(cmds, field.Init())
	}

	cmd := g.fields[g.current].Focus()
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

// setCurrent sets the current field.
func (g *Group) setCurrent(current int) tea.Cmd {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	cmd = g.fields[g.current].Blur()
	cmds = append(cmds, cmd)

	g.current = ordered.Clamp(current, 0, len(g.fields)-1)

	cmd = g.fields[g.current].Focus()
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

// Update updates the group.
func (g *Group) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	m, cmd := g.fields[g.current].Update(msg)
	g.fields[g.current] = m.(Field)

	cmds = append(cmds, cmd)

	switch msg.(type) {
	case nextFieldMsg:
		current := g.current
		cmd = g.setCurrent(current + 1)

		if current == len(g.fields)-1 {
			cmds = append(cmds, nextGroup)
			break
		}

		cmds = append(cmds, cmd)

	case prevFieldMsg:
		current := g.current
		cmd = g.setCurrent(current - 1)

		if current == 0 {
			cmds = append(cmds, prevGroup)
			break
		}

		cmds = append(cmds, cmd)
	}

	return g, tea.Batch(cmds...)
}

// View renders the group.
func (g *Group) View() string {
	var s strings.Builder

	gap := g.theme.FieldSeparator.String()
	if gap == "" {
		gap = "\n\n"
	}

	for i, field := range g.fields {
		s.WriteString(field.View())
		if i < len(g.fields)-1 {
			s.WriteString(gap)
		}
	}

	if g.showHelp {
		s.WriteString(gap)
		s.WriteString(g.theme.Focused.Help.Render(g.help.ShortHelpView(g.fields[g.current].KeyBinds())))
		s.WriteString("\n")
	}

	for _, err := range g.Errors() {
		s.WriteString("\n")
		s.WriteString(g.theme.Focused.ErrorMessage.Render(err.Error()))
	}

	if len(g.Errors()) == 0 {
		// If there are no errors add a gap so that the appearance of an
		// error message doesn't cause the layout to shift.
		//
		// XXX: Mutli-line errors will still cause a shift. How do we handle
		// this?
		s.WriteString("\n")
	}

	return s.String()
}
