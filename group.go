package huh

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Group is a collection of fields displayed together.
type Group struct {
	fields []Field

	title       string
	description string
	current     int
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

// Errors returns the group's errors.
func (g *Group) Errors() []error {
	var errors []error
	for _, field := range g.fields {
		if err := field.Error(); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
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

	if g.title != "" {
		s.WriteString(g.title)
		s.WriteString("\n")
	}

	if g.description != "" {
		s.WriteString(g.description)
		s.WriteString("\n")
	}

	for _, field := range g.fields {
		s.WriteString(field.View())
		s.WriteString("\n")
	}

	if len(g.Errors()) > 0 {
		for _, err := range g.Errors() {
			s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(" * " + err.Error()))
			s.WriteString("\n")
		}
	}

	return s.String()
}
