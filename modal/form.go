package modal

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
)

// formContent adapts a [huh.Form] to the [Content] interface. The form's
// own State drives the modal's resolution: [huh.StateCompleted] → confirmed,
// [huh.StateAborted] → cancelled.
type formContent struct {
	form *huh.Form
}

func (f *formContent) Init() tea.Cmd { return f.form.Init() }

func (f *formContent) Update(msg tea.Msg) (Content, tea.Cmd) {
	model, cmd := f.form.Update(msg)
	if updated, ok := model.(*huh.Form); ok {
		f.form = updated
	}
	return f, cmd
}

func (f *formContent) View() string { return f.form.View() }

func (f *formContent) Done() (done, confirmed bool, value any) {
	switch f.form.State {
	case huh.StateCompleted:
		return true, true, f.form
	case huh.StateAborted:
		return true, false, nil
	default:
		return false, false, nil
	}
}

// NewForm constructs a Modal that hosts a huh form. When the form completes
// the modal resolves with Confirmed=true and Value set to the [*huh.Form]
// (read field values via [huh.Form.GetString] etc.). When the user aborts
// the form (esc by default) the modal resolves with Confirmed=false.
func NewForm(id string, form *huh.Form, opts ...Option) *Modal {
	return New(id, &formContent{form: form}, opts...)
}
