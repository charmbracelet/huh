package spinner

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
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
	ctx        context.Context
	accessible bool
	output     *termenv.Output
	title      string
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

// Context sets the context of the spinner.
func (s *Spinner) Context(ctx context.Context) *Spinner {
	s.ctx = ctx
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

// Accessible sets the spinner to be static.
func (s *Spinner) Accessible(accessible bool) *Spinner {
	s.accessible = accessible
	return s
}

// New creates a new spinner.
func New() *Spinner {
	s := spinner.New()
	s.Spinner = spinner.Dot

	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#F780E2"))

	return &Spinner{
		action:     func() { time.Sleep(time.Second) },
		spinner:    s,
		title:      "Loading...",
		titleStyle: lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}),
		output:     termenv.NewOutput(os.Stdout),
		ctx:        nil,
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
	if s.accessible {
		return s.runAccessible()
	}

	p := tea.NewProgram(s, tea.WithContext(s.ctx), tea.WithOutput(os.Stderr))

	if s.ctx == nil {
		go func() {
			s.action()
			p.Quit()
		}()
	}

	_, _ = p.Run()

	return nil
}

// runAccessible runs the spinner in an accessible mode (statically).
func (s *Spinner) runAccessible() error {
	s.output.HideCursor()
	frame := s.spinner.Style.Render("...")
	title := s.titleStyle.Render(strings.TrimSuffix(s.title, "..."))
	fmt.Println(title + frame)

	if s.ctx == nil {
		s.action()
		s.output.ShowCursor()
		s.output.CursorBack(len(frame) + len(title))
		return nil
	}

	actionDone := make(chan struct{})

	go func() {
		s.action()
		actionDone <- struct{}{}
	}()

	for {
		select {
		case <-s.ctx.Done():
			s.output.ShowCursor()
			s.output.CursorBack(len(frame) + len(title))
			return s.ctx.Err()
		case <-actionDone:
			s.output.ShowCursor()
			s.output.CursorBack(len(frame) + len(title))
			return nil
		}
	}
}
