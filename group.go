package huh

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Group is a collection of fields displayed together.
type Group struct {
	fields []Field

	label       string
	description string
	current     int
	errors      []error
}

// NewGroup creates a new group with the given fields.
func NewGroup(fields ...Field) *Group {
	return &Group{
		fields:  fields,
		current: 0,
	}
}

// Label sets the group's label.
func (g *Group) Label(label string) *Group {
	g.label = label
	return g
}

// Description sets the group's description.
func (g *Group) Description(description string) *Group {
	g.description = description
	return g
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
	for i, field := range g.fields {
		field.Init()
		if i == 0 {
			field.Focus()
		}
	}
	return nil
}

// setCurrent sets the current field.
func (g *Group) setCurrent(current int) {
	g.fields[g.current].Blur()
	if current < 0 {
		current = 0
	}
	if current > len(g.fields)-1 {
		current = len(g.fields) - 1
	}
	g.current = current
	g.fields[g.current].Focus()
}

// Update updates the group.
func (g *Group) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	m, cmd := g.fields[g.current].Update(msg)
	g.fields[g.current] = m.(Field)

	cmds = append(cmds, cmd)

	switch msg.(type) {
	case nextFieldMsg:
		if g.current == len(g.fields)-1 {
			cmds = append(cmds, nextGroup)
			break
		}

		g.setCurrent(g.current + 1)

	case prevFieldMsg:
		if g.current == 0 {
			cmds = append(cmds, prevGroup)
			break
		}

		g.setCurrent(g.current - 1)
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

	return s.String()
}
