package spinner

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Spinner represents a loading spinner.
// To get started simply create a new spinner and call `Run`.
//
//	s := spinner.New()
//	s.Run()
//
// ⣾  Loading...
type Spinner struct {
	spinner    spinner.Model
	action     func()
	title      string
	style      lipgloss.Style
	titleStyle lipgloss.Style
}

type Type spinner.Spinner

var (
	Line      = Type(spinner.Line)
	Dots      = Type(spinner.Dot)
	MiniDot   = Type(spinner.MiniDot)
	Jump      = Type(spinner.Jump)
	Points    = Type(spinner.Points)
	Pulse     = Type(spinner.Pulse)
	Globe     = Type(spinner.Globe)
	Moon      = Type(spinner.Moon)
	Monkey    = Type(spinner.Monkey)
	Meter     = Type(spinner.Meter)
	Hamburger = Type(spinner.Hamburger)
	Ellipsis  = Type(spinner.Ellipsis)
)

// Type sets the type of the spinner.
func (s *Spinner) Type(t Type) *Spinner {
	s.spinner.Spinner = spinner.Spinner(t)
	return s
}

// Title sets the title of the spinner.
func (s *Spinner) Title(title string) *Spinner {
	s.title = title
	return s
}

// Action sets the action of the spinner.
func (s *Spinner) Action(action func()) *Spinner {
	s.action = action
	return s
}

// Style sets the style of the spinner.
func (s *Spinner) Style(style lipgloss.Style) *Spinner {
	s.spinner.Style = style
	return s
}

// TitleStyle sets the title style of the spinner.
func (s *Spinner) TitleStyle(style lipgloss.Style) *Spinner {
	s.titleStyle = style
	return s
}

// New creates a new spinner.
func New() *Spinner {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return &Spinner{
		action:     func() { time.Sleep(time.Second) },
		spinner:    s,
		title:      "Loading...",
		titleStyle: lipgloss.NewStyle(),
	}
}

// Init initializes the spinner.
func (s *Spinner) Init() tea.Cmd {
	return s.spinner.Tick
}

// Update updates the spinner.
func (s *Spinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return s, tea.Quit
		}
	}

	var cmd tea.Cmd
	s.spinner, cmd = s.spinner.Update(msg)
	return s, cmd
}

// View returns the spinner view.
func (s *Spinner) View() string {
	var title string
	if s.title != "" {
		title = s.titleStyle.Render(s.title) + " "
	}
	return s.spinner.View() + title
}

// Run runs the spinner.
func (s *Spinner) Run() error {
	p := tea.NewProgram(s)

	go func() {
		s.action()
		p.Quit()
	}()

	_, _ = p.Run()

	return nil
}
