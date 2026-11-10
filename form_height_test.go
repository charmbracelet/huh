package huh

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// A form whose field has a dynamic description must resize when the
// description changes, otherwise the extra lines are cut off.
func TestFormHeightFollowsDynamicDescription(t *testing.T) {
	showHelp := false
	field := NewInput().
		Title("Enter a string").
		DescriptionFunc(func() string {
			desc := "This is a test prompt"
			if showHelp {
				desc += "\n\nEnter a value at least 5 characters long"
			}
			return desc
		}, &showHelp)

	f := NewForm(NewGroup(field)).WithWidth(40)

	m, _ := f.Update(tea.WindowSizeMsg{Width: 40, Height: 30})
	f = m.(*Form)

	before := lipgloss.Height(f.View())

	// Change the binding so that the description grows.
	showHelp = true

	m, cmd := f.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})
	f = m.(*Form)

	// Deliver the messages the form asked for, such as the updated
	// description.
	deliver := func(cmd tea.Cmd) {
		if cmd == nil {
			return
		}
		msg := cmd()
		if batch, ok := msg.(tea.BatchMsg); ok {
			for _, c := range batch {
				if c == nil {
					continue
				}
				if sub := c(); sub != nil {
					m, _ = f.Update(sub)
					f = m.(*Form)
				}
			}
			return
		}
		if msg != nil {
			m, _ = f.Update(msg)
			f = m.(*Form)
		}
	}
	deliver(cmd)

	after := lipgloss.Height(f.View())
	if after <= before {
		t.Fatalf("form height did not grow after the description changed: before=%d after=%d\n%s", before, after, f.View())
	}
}
