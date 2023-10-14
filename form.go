package huh

import (
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
)

// Form represents a Huh? form.
// It is a collection of groups and controls navigation between pages.
type Form struct {
	groups    []*Group
	paginator paginator.Model
	quitting  bool
}

// NewForm creates a new form with the given groups.
func NewForm(groups ...*Group) *Form {
	p := paginator.New()
	p.SetTotalPages(len(groups))

	return &Form{
		groups:    groups,
		paginator: p,
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
			f.quitting = true
			return f, tea.Quit
		}

	case nextGroupMsg:
		if f.paginator.OnLastPage() {
			f.quitting = true
			return f, tea.Quit
		}
		f.paginator.NextPage()

	case prevGroupMsg:
		f.paginator.PrevPage()
	}

	m, cmd := f.groups[f.paginator.Page].Update(msg)
	f.groups[f.paginator.Page] = m.(*Group)

	return f, cmd
}

// View renders the form.
func (f *Form) View() string {
	if f.quitting {
		return ""
	}
	return f.groups[f.paginator.Page].View()
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
