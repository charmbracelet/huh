package huh

import (
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
)

// Form represents a Huh? form.
// It is a collection of groups and controls navigation between pages.
type Form struct {
	groups     []*Group
	paginator  paginator.Model
	accessible bool
	quitting   bool
	theme      *Theme
}

// NewForm creates a new form with the given groups.
func NewForm(groups ...*Group) *Form {
	p := paginator.New()
	p.SetTotalPages(len(groups))

	return &Form{
		groups:    groups,
		paginator: p,
		theme:     NewBaseTheme(),
	}
}

// Field is a form field.
type Field interface {
	// Bubble Tea Model
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string

	// Bubble Tea Events
	Blur() tea.Cmd
	Focus() tea.Cmd

	// Accessible Prompt (non-redraw)
	Run()

	// Theme sets the theme on a field.
	Theme(*Theme) Field
}

type nextGroupMsg struct{}
type prevGroupMsg struct{}

func nextGroup() tea.Msg {
	return nextGroupMsg{}
}

func prevGroup() tea.Msg {
	return prevGroupMsg{}
}

// Accessible sets the form to run in accessible mode to avoid redrawing the
// views which makes it easier for screen readers to read and describe the form.
//
// This avoids using the Bubble Tea renderer and instead simply uses basic
// terminal prompting to gather input which degrades the user experience but
// provides accessibility.
func (f *Form) Accessible(b bool) *Form {
	f.accessible = b
	return f
}

// Init initializes the form.
func (f *Form) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, group := range f.groups {
		cmds = append(cmds, group.Init())
	}
	return tea.Batch(cmds...)
}

// Update updates the form.
func (f *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
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
	// Make theme acessible to groups and fields
	for _, group := range f.groups {
		group.Theme(f.theme)
		for _, field := range group.fields {
			field.Theme(f.theme)
		}
	}

	if f.accessible {
		for _, group := range f.groups {
			for _, field := range group.fields {
				field.Init()
				field.Focus()
				field.Run()
			}
		}
		return nil
	}

	p := tea.NewProgram(f)
	_, err := p.Run()
	if err != nil {
		return err
	}
	return nil
}
