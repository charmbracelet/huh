package huh

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/ordered"
)

// Form represents a Huh? form.
// It is a collection of groups and controls navigation between pages.
type Form struct {
	groups []*Group
	page   int
}

// NewForm creates a new form with the given groups.
func NewForm(groups ...*Group) *Form {
	return &Form{
		groups: groups,
		page:   0,
	}
}

type nextGroupMsg struct{}
type prevGroupMsg struct{}

func nextGroup() tea.Msg {
	return nextGroupMsg{}
}

func prevGroup() tea.Msg {
	return prevGroupMsg{}
}

// Init initializes the form.
func (f *Form) Init() tea.Cmd {
	for _, group := range f.groups {
		group.Init()
	}
	return nil
}

// Update updates the form.
func (f *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return f, tea.Quit
		}

	case nextGroupMsg:
		if f.page == len(f.groups)-1 {
			return f, tea.Quit
		}
		f.page = ordered.Min(f.page+1, len(f.groups)-1)

	case prevGroupMsg:
		f.page = ordered.Max(f.page-1, 0)
	}

	m, cmd := f.groups[f.page].Update(msg)
	f.groups[f.page] = m.(*Group)

	return f, cmd
}

// View renders the form.
func (f *Form) View() string {
	return f.groups[f.page].View()
}

// Run runs the form.
func (f *Form) Run() error {
	p := tea.NewProgram(f)
	_, err := p.Run()
	if err != nil {
		return err
	}
	return nil
}
