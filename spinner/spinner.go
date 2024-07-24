package spinner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

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
	action     func(ctx context.Context) error
	ctx        context.Context
	accessible bool
	output     *termenv.Output
	title      string
	titleStyle lipgloss.Style

	err error
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
	s.action = func(ctx context.Context) error {
		action()
		return nil
	}
	return s
}

// ActionErr sets the action of the spinner.
func (s *Spinner) ActionErr(action func(ctx context.Context) error) *Spinner {
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
		spinner:    s,
		title:      "Loading...",
		titleStyle: lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}),
		output:     termenv.NewOutput(os.Stdout),
	}
}

// Init initializes the spinner.
func (s *Spinner) Init() tea.Cmd {
	return tea.Batch(s.spinner.Tick, func() tea.Msg {
		if s.action != nil {
			return doneMsg{err: s.action(s.ctx)}
		}
		return nil
	})
}

// Update updates the spinner.
func (s *Spinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case doneMsg:
		s.err = msg.err
		return s, tea.Quit
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
	hasCtx := s.ctx != nil

	if hasCtx && s.ctx.Err() != nil {
		if errors.Is(s.ctx.Err(), context.Canceled) {
			return nil
		}
		return s.ctx.Err()
	}

	// sets a dummy action if the spinner does not have a context nor an action.
	if !hasCtx && s.action == nil {
		// there's nothing to do!
		return nil
	}

	if s.accessible {
		return s.runAccessible()
	}

	m, err := tea.NewProgram(s, tea.WithContext(s.ctx), tea.WithOutput(os.Stderr)).Run()
	mm := m.(*Spinner)
	if mm.err != nil {
		return mm.err
	}
	return err
}

// runAccessible runs the spinner in an accessible mode (statically).
func (s *Spinner) runAccessible() error {
	s.output.HideCursor()
	frame := s.spinner.Style.Render("...")
	title := s.titleStyle.Render(strings.TrimSuffix(s.title, "..."))
	fmt.Println(title + frame)

	if s.ctx == nil {
		err := s.action(context.Background())
		s.output.ShowCursor()
		s.output.CursorBack(len(frame) + len(title))
		return err
	}

	actionDone := make(chan error)
	if s.action != nil {
		go func() {
			actionDone <- s.action(s.ctx)
		}()
	}

	for {
		select {
		case <-s.ctx.Done():
			s.output.ShowCursor()
			s.output.CursorBack(len(frame) + len(title))
			if errors.Is(s.ctx.Err(), context.Canceled) {
				return nil
			}
			return s.ctx.Err()
		case err := <-actionDone:
			s.output.ShowCursor()
			s.output.CursorBack(len(frame) + len(title))
			return err
		}
	}
}

type doneMsg struct {
	err error
}
