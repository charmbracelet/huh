package huh

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Group is a collection of fields displayed together.
type Group struct {
	fields []Field

	title       string
	description string
	current     int
	theme       *Theme
}

// NewGroup creates a new group with the given fields.
func NewGroup(fields ...Field) *Group {
	return &Group{
		fields:  fields,
		current: 0,
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

// Theme sets the theme on a group.
func (g *Group) Theme(t *Theme) *Group {
	g.theme = t
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

	for _, field := range g.fields {
		s.WriteString(field.View())
		s.WriteString("\n")
	}

	for _, err := range g.Errors() {
		s.WriteString(g.theme.Focused.Error.Render(err.Error()))
		s.WriteString("\n")
	}

	return s.String()
}
