package spinner

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestNewSpinner(t *testing.T) {
	s := New()
	if s.title != "Loading..." {
		t.Errorf("Expected default title 'Loading...', got '%s'", s.title)
	}
	if !reflect.DeepEqual(s.spinner.Spinner, spinner.Dot) {
		t.Errorf("Expected default spinner type to be Dot, got %v", s.spinner.Spinner)
	}
}

func TestSpinnerType(t *testing.T) {
	s := New().Type(Dots)
	if !reflect.DeepEqual(s.spinner.Spinner, spinner.Dot) {
		t.Errorf("Expected spinner type to be Dot, got %v", s.spinner.Spinner)
	}
}

func TestSpinnerDifferentTypes(t *testing.T) {
	s := New().Type(Line)
	if !reflect.DeepEqual(s.spinner.Spinner, spinner.Line) {
		t.Errorf("Expected spinner type to be Line, got %v", s.spinner.Spinner)
	}
}

func TestSpinnerView(t *testing.T) {
	s := New().Title("Test")
	view := s.View()

	if !strings.Contains(view, "Test") {
		t.Errorf("Expected view to contain title 'Test', got '%s'", view)
	}
}

func TestSpinnerContextCancellation(t *testing.T) {
	exercise(t, func() *Spinner {
		ctx, cancel := context.WithCancel(context.Background())
		s := New().Context(ctx)
		cancel() // Cancel before running
		return s
	}, requireContextCanceled)
}

func TestSpinnerContextCancellationWhileRunning(t *testing.T) {
	exercise(t, func() *Spinner {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()
		return New().Accessible(true).Context(ctx)
	}, requireContextCanceled)
}

func TestSpinnerStyleMethods(t *testing.T) {
	s := New()
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("red"))
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("blue"))

	s.Style(style)
	s.TitleStyle(titleStyle)

	if !reflect.DeepEqual(s.spinner.Style, style) {
		t.Errorf("Style was not set correctly")
	}

	if !reflect.DeepEqual(s.titleStyle, titleStyle) {
		t.Errorf("TitleStyle was not set correctly")
	}
}

func TestSpinnerInit(t *testing.T) {
	s := New()
	cmd := s.Init()

	if cmd == nil {
		t.Errorf("Init did not return a valid command")
	}
}

func TestSpinnerUpdate(t *testing.T) {
	s := New()
	cmd := s.Init()
	if cmd == nil {
		t.Errorf("Init did not return a valid command")
	}

	model, cmd := s.Update(spinner.TickMsg{})
	if reflect.TypeOf(model) != reflect.TypeOf(&Spinner{}) {
		t.Errorf("Update did not return correct model type")
	}

	if cmd == nil {
		t.Errorf("Update should return a non-nil command in this scenario")
	}

	// Simulate key press
	_, cmd = s.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Errorf("Update did not handle key press correctly")
	}
}

func TestSpinnerSimple(t *testing.T) {
	exercise(t, func() *Spinner {
		return New().Action(func() {})
	}, requireNoError)
}

func TestSpinnerWithContextAndAction(t *testing.T) {
	exercise(t, func() *Spinner {
		ctx := context.Background()
		return New().Context(ctx).Action(func() {})
	}, requireNoError)
}

func TestSpinnerWithActionError(t *testing.T) {
	fake := errors.New("fake")
	exercise(t, func() *Spinner {
		return New().ActionWithErr(func(context.Context) error { return fake })
	}, requireErrorIs(fake))
}

func exercise(t *testing.T, factory func() *Spinner, checker func(tb testing.TB, err error)) {
	t.Helper()
	t.Run("accessible", func(t *testing.T) {
		err := factory().Accessible(true).Run()
		checker(t, err)
	})
	t.Run("regular", func(t *testing.T) {
		err := factory().Accessible(false).Run()
		checker(t, err)
	})
}

func requireNoError(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Errorf("expected no error, got %v", err)
	}
}

func requireErrorIs(target error) func(tb testing.TB, err error) {
	return func(tb testing.TB, err error) {
		tb.Helper()
		if !errors.Is(err, target) {
			tb.Errorf("expected error to be %v, got %v", target, err)
		}
	}
}

func requireContextCanceled(tb testing.TB, err error) {
	tb.Helper()
	if !errors.Is(err, context.Canceled) && !errors.Is(err, tea.ErrProgramKilled) {
		tb.Errorf("expected to get a context canceled error, got %v", err)
	}
}
