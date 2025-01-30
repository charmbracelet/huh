package spinner

import (
	"context"
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
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
	ctx, cancel := context.WithCancel(context.Background())

	s := New().Context(ctx)
	cancel() // Cancel before running

	err := s.Run()
	if err != nil {
		t.Errorf("Run() returned an error after context cancellation: %v", err)
	}
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

func TestAccessibleSpinner(t *testing.T) {
	s := New().Accessible(true)
	err := s.Run()
	if err != nil {
		t.Errorf("Run() in accessible mode returned an error: %v", err)
	}
}

func TestSpinnerOutput(t *testing.T) {
	tests := []struct {
		name       string
		wantStdout bool
		wantStderr bool
	}{
		{
			name:       "stdout",
			wantStdout: true,
			wantStderr: false,
		},
		{
			name:       "stderr",
			wantStdout: false,
			wantStderr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			const title = "Test Output"

			// Save original stderr and stdout
			oldStderr := os.Stderr
			oldStdout := os.Stdout

			// Create pipes for stderr and stdout
			stderrReader, stderrWriter, _ := os.Pipe()
			stdoutReader, stdoutWriter, _ := os.Pipe()

			// Set global stderr and stdout to our pipes
			os.Stderr = stderrWriter
			os.Stdout = stdoutWriter

			// Create a spinner and set its output
			s := New().Title(title).Accessible(true)
			if tc.wantStderr {
				s.Output(termenv.NewOutput(os.Stderr))
			}
			if tc.wantStdout {
				s.Output(termenv.NewOutput(os.Stdout))
			}
			s.action = func() { time.Sleep(100 * time.Millisecond) }
			if err := s.Run(); err != nil {
				t.Errorf("Spinner.Run() returned an error: %v", err)
			}

			// Restore original stderr and stdout
			os.Stderr = oldStderr
			os.Stdout = oldStdout

			// Close the pipes
			if err := errors.Join(stderrWriter.Close(), stdoutWriter.Close()); err != nil {
				t.Errorf("Failed to close pipes: %v", err)
			}

			// Read from the pipes
			stderrOutput, stderrErr := io.ReadAll(stderrReader)
			stdoutOutput, stdoutErr := io.ReadAll(stdoutReader)
			if err := errors.Join(stderrErr, stdoutErr); err != nil {
				t.Errorf("Failed to read from pipes: %v", err)
			}

			// Check the output
			if tc.wantStderr {
				if !strings.Contains(string(stderrOutput), title) {
					t.Errorf("Stderr got %q, but wanted %q", stderrOutput, title)
				}
				if len(stdoutOutput) > 0 {
					t.Errorf("Expected no output on stdout, but got %q", stdoutOutput)
				}
			}
			if tc.wantStdout {
				if !strings.Contains(string(stdoutOutput), title) {
					t.Errorf("Stdout got %q, but wanted %q", stdoutOutput, title)
				}
				if len(stderrOutput) > 0 {
					t.Errorf("Expected no output on stderr, but got %q", stderrOutput)
				}
			}
		})
	}
}
