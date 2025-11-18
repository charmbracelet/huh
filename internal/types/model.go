package types

import tea "charm.land/bubbletea/v2"

type Model interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Model, tea.Cmd)
	View() string
}

type ViewModel struct {
	Model
	ViewHook func(view tea.View) tea.View
}

func (w ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m, cmd := w.Model.Update(msg)
	return ViewModel{Model: m}, cmd
}

func (w ViewModel) View() tea.View {
	var view tea.View
	if w.ViewHook != nil {
		view = w.ViewHook(view)
	}
	view.SetContent(w.Model.View())
	return view
}
