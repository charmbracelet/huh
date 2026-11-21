package huh

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func escKeyPress() tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: tea.KeyEscape})
}

func testSelectWithTwoOptions() *Select[string] {
	var value string
	return NewSelect[string]().
		Options(NewOptions("one", "two")...).
		Value(&value)
}

func pressEscape(t *testing.T, s *Select[string]) {
	t.Helper()
	s.Update(escKeyPress())
}

func assertHovered(t *testing.T, s *Select[string], want string) {
	t.Helper()
	got, ok := s.Hovered()
	if !ok || got != want {
		t.Fatalf("hovered option = %q (present=%v), want %q", got, ok, want)
	}
}

// TestSelectInitialFilteringKeyMap verifies that a Select constructed with
// Filtering(true) starts in filter mode regardless of whether its key map is
// assigned before or after that option is set. The form always overwrites the
// key map while constructing the fields, so the binding state must be
// reconciled after the assignment rather than only by Filtering itself
// (issue #854).
func TestSelectInitialFilteringKeyMap(t *testing.T) {
	cases := []struct {
		name  string
		build func() *Select[string]
	}{
		{
			name: "filtering set before the key map",
			build: func() *Select[string] {
				s := testSelectWithTwoOptions().Filtering(true)
				s.WithKeyMap(NewDefaultKeyMap())
				return s
			},
		},
		{
			name: "filtering set after the key map",
			build: func() *Select[string] {
				s := testSelectWithTwoOptions()
				s.WithKeyMap(NewDefaultKeyMap())
				s.Filtering(true)
				return s
			},
		},
		{
			name: "form applies the key map during construction",
			build: func() *Select[string] {
				s := testSelectWithTwoOptions().Filtering(true)
				f := NewForm(NewGroup(s))
				f.Update(f.Init())
				return s
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := tc.build()

			if !s.GetFiltering() {
				t.Fatal("the select should start in filter mode")
			}
			if !s.keymap.SetFilter.Enabled() {
				t.Error("SetFilter (esc) should be enabled while filter mode is active")
			}
			if s.keymap.Filter.Enabled() {
				t.Error("Filter (/) should be disabled while filter mode is active")
			}
			if s.keymap.ClearFilter.Enabled() {
				t.Error("ClearFilter should stay disabled while the filter is empty")
			}

			// Esc leaves the initial filter mode, exactly as it does after
			// pressing "/" and then Esc in a normally constructed select.
			pressEscape(t, s)
			if s.GetFiltering() {
				t.Fatal("the select should leave filter mode after Esc")
			}
			if s.keymap.SetFilter.Enabled() {
				t.Error("SetFilter (esc) should be disabled outside filter mode")
			}
			if !s.keymap.Filter.Enabled() {
				t.Error("Filter (/) should be enabled outside filter mode")
			}

			// Navigation must work again once filter mode is left.
			m, _ := s.Update(keypress('j'))
			s = m.(*Select[string])
			assertHovered(t, s, "two")
		})
	}
}

// TestSelectEscWithoutFilterMode ensures Esc keeps its existing behavior for
// selects that do not start in filter mode.
func TestSelectEscWithoutFilterMode(t *testing.T) {
	s := testSelectWithTwoOptions()
	s.WithKeyMap(NewDefaultKeyMap())

	pressEscape(t, s)
	if s.GetFiltering() {
		t.Fatal("filter mode should not be active")
	}

	m, _ := s.Update(codeKeypress(tea.KeyDown))
	s = m.(*Select[string])
	assertHovered(t, s, "two")
}

// TestSelectInitialFilteringHelp renders the form help line for a select that
// starts in filter mode. The help must describe Esc as the way to leave the
// initial filter mode and must not advertise "/" while filter mode is active.
func TestSelectInitialFilteringHelp(t *testing.T) {
	var value string
	s := NewSelect[string]().
		Options(NewOptions("one", "two")...).
		Value(&value).
		Filtering(true)
	f := NewForm(NewGroup(s))
	f.Update(f.Init())

	view := viewModel(f)
	if !strings.Contains(view, "esc set filter") {
		t.Errorf("initial help should offer Esc to leave filter mode, got:\n%s", view)
	}
	if strings.Contains(view, "/ filter") {
		t.Errorf("initial help should not offer / to enter filter mode, got:\n%s", view)
	}

	m, _ := f.Update(escKeyPress())
	f = m.(*Form)
	view = viewModel(f)
	if !strings.Contains(view, "/ filter") {
		t.Errorf("help should offer / to enter filter mode after Esc, got:\n%s", view)
	}
	if strings.Contains(view, "esc set filter") {
		t.Errorf("help should not offer Esc to leave filter mode after Esc, got:\n%s", view)
	}

	m, _ = f.Update(keypress('j'))
	_ = m
	assertHovered(t, s, "two")
}
