package huh

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

// Group is a collection of fields displayed together.
type Group struct {
	fields []Field

	title       string
	description string
	current     int

	showHelp bool
	help     help.Model

	theme  *Theme
	keymap *KeyMap
}

// NewGroup creates a new group with the given fields.
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
		err := field.Error()
		if err != nil {
			errs = append(errs, field.Error())
		}
	}
	return errs
}

type nextFieldMsg struct{}
type prevFieldMsg struct{}

func nextField() tea.Msg {
	return nextFieldMsg{}
}

func prevField() tea.Msg {
	return prevFieldMsg{}
}

// Init initializes the group.
func (g *Group) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, field := range g.fields {
		cmds = append(cmds, field.Init())
	}
	cmds = append(cmds, g.fields[g.current].Focus())
	return tea.Batch(cmds...)
}

// setCurrent sets the current field.
func (g *Group) setCurrent(current int) tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, g.fields[g.current].Blur())
	if current < 0 {
		current = 0
	}
	if current > len(g.fields)-1 {
		current = len(g.fields) - 1
	}
	g.current = current
	cmds = append(cmds, g.fields[g.current].Focus())
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
