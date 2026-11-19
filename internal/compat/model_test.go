package compat

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

type textModel string

func (textModel) Init() tea.Cmd                     { return nil }
func (m textModel) Update(tea.Msg) (Model, tea.Cmd) { return m, nil }
func (m textModel) View() string                    { return string(m) }

func TestViewModelViewHook(t *testing.T) {
	const content = "rendered model"
	const replacement = "custom view"

	model := ViewModel{
		Model: textModel(content),
		ViewHook: func(view tea.View) tea.View {
			if view.Content != content {
				t.Errorf("hook content = %q, want %q", view.Content, content)
			}
			view.SetContent(replacement)
			view.AltScreen = true
			return view
		},
	}

	view := model.View()
	if view.Content != replacement {
		t.Errorf("returned content = %q, want %q", view.Content, replacement)
	}
	if !view.AltScreen {
		t.Error("returned view lost the hook's AltScreen setting")
	}

	model.ViewHook = nil
	if got := model.View().Content; got != content {
		t.Errorf("content without a hook = %q, want %q", got, content)
	}
}
