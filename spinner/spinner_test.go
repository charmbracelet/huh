package spinner

import (
	"context"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
			time.Sleep(250 * time.Millisecond)
			cancel()
		}()
		return New().Context(ctx)
	}, requireContextCanceled)
}

func TestSpinnerStyleMethods(t *testing.T) {
	s := New()

	theme := ThemeFunc(func(bool) *Styles {
		return &Styles{
			Spinner: lipgloss.NewStyle().Foreground(lipgloss.Color("red")),
			Title:   lipgloss.NewStyle().Foreground(lipgloss.Color("blue")),
		}
	})

	s.WithTheme(theme).View()
	styles := s.theme.Theme(true)
	if !reflect.DeepEqual(s.spinner.Style, styles.Spinner) {
		t.Errorf("Style was not set correctly")
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
	_, cmd = s.Update(tea.KeyPressMsg(tea.Key{Mod: tea.ModCtrl, Code: 'c'}))
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
		err := factory().
			WithAccessible(true).
			WithOutput(io.Discard).
			WithInput(nilReader{}).
			Run()
		checker(t, err)
	})
	t.Run("regular", func(t *testing.T) {
		err := factory().
			WithAccessible(false).
			WithOutput(io.Discard).
			WithInput(nilReader{}).
			Run()
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
	switch {
	case errors.Is(err, context.Canceled):
	case errors.Is(err, tea.ErrProgramKilled):
	default:
		tb.Errorf("expected to get a context canceled error, got %v", err)
	}
}

type nilReader struct{}

// Read implements io.Reader.
func (nilReader) Read([]byte) (int, error) { return 0, nil }
