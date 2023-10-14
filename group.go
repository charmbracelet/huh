package huh

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/ordered"
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
	for _, field := range g.fields {
		field.Init()
	}
	return nil
}

// Update updates the group.
func (g *Group) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	g.fields[g.current], cmd = g.fields[g.current].Update(msg)

	cmds = append(cmds, cmd)

	switch msg.(type) {
	case nextFieldMsg:
		if g.current == len(g.fields)-1 {
			cmds = append(cmds, nextGroup)
			break
		}
		g.current = ordered.Min(g.current+1, len(g.fields)-1)
	case prevFieldMsg:
		if g.current == 0 {
			cmds = append(cmds, prevGroup)
			break
		}
		g.current = ordered.Max(g.current-1, 0)
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
