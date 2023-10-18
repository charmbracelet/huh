package huh

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// SelectStyle is the style of the select field.
type SelectStyle struct {
	Title      lipgloss.Style
	Cursor     lipgloss.Style
	Selected   lipgloss.Style
	Unselected lipgloss.Style
}

// DefaultSelectStyles returns the default focused style of the select field.
func DefaultSelectStyles() (SelectStyle, SelectStyle) {
	focused := SelectStyle{
		Title:      lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		Cursor:     lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		Selected:   lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
		Unselected: lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
	}
	blurred := SelectStyle{
		Title:      lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Cursor:     lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Selected:   lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Unselected: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
	}
	return focused, blurred
}

// MultiSelectStyle is the style of the multi-select field.
type MultiSelectStyle struct {
	Title            lipgloss.Style
	Cursor           lipgloss.Style
	Selected         lipgloss.Style
	Unselected       lipgloss.Style
	SelectedPrefix   lipgloss.Style
	UnselectedPrefix lipgloss.Style
}

// DefaultMultiSelectStyles returns the default focused style of the multi-select field.
func DefaultMultiSelectStyles() (MultiSelectStyle, MultiSelectStyle) {
	focused := MultiSelectStyle{
		Title:            lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		Cursor:           lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		Selected:         lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
		Unselected:       lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
		SelectedPrefix:   lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
		UnselectedPrefix: lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
	}
	blurred := MultiSelectStyle{
		Title:            lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Cursor:           lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Selected:         lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Unselected:       lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		SelectedPrefix:   lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		UnselectedPrefix: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
	}
	return focused, blurred
}

// TextareaStyle is the style of the textarea field.
type TextStyle struct {
	Title lipgloss.Style
	Help  lipgloss.Style
	textarea.Style
}

// DefaultTextStyles returns the default focused style of the text field.
func DefaultTextStyles() (TextStyle, TextStyle) {
	f, b := textarea.DefaultStyles()
	f.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	b.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	f.CursorLine = lipgloss.NewStyle()
	b.CursorLine = lipgloss.NewStyle()

	focused := TextStyle{
		Title: lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		Help:  lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Style: f,
	}
	blurred := TextStyle{
		Title: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Help:  lipgloss.NewStyle().Foreground(lipgloss.Color("0")),
		Style: b,
	}

	return focused, blurred
}

// InputStyle is the style of the input field.
type InputStyle struct {
	Title       lipgloss.Style
	Prompt      lipgloss.Style
	Text        lipgloss.Style
	Placeholder lipgloss.Style
}

// DefaultInputStyles returns the default focused style of the input field.
func DefaultInputStyles() (InputStyle, InputStyle) {
	focused := InputStyle{
		Title:       lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		Prompt:      lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
		Text:        lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
		Placeholder: lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
	}
	blurred := InputStyle{
		Title:       lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Prompt:      lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Text:        lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Placeholder: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
	}
	return focused, blurred
}

// ConfirmStyle is the style of the confirm field.
type ConfirmStyle struct {
	Title      lipgloss.Style
	Selected   lipgloss.Style
	Unselected lipgloss.Style
}

// DefaultConfirmStyles returns the default focused style of the confirm field.
func DefaultConfirmStyles() (ConfirmStyle, ConfirmStyle) {
	focused := ConfirmStyle{
		Title:      lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		Selected:   lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Background(lipgloss.Color("4")).Padding(0, 2).Margin(1),
		Unselected: lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Background(lipgloss.Color("0")).Padding(0, 2).Margin(1),
	}
	blurred := ConfirmStyle{
		Title:      lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Selected:   lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Background(lipgloss.Color("0")).Padding(0, 2).Margin(1),
		Unselected: lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Background(lipgloss.Color("0")).Padding(0, 2).Margin(1),
	}
	return focused, blurred
}

// NoteStyle is the style of the Note field.
type NoteStyle struct {
	Title lipgloss.Style
	Body  lipgloss.Style
}

// DefaultNoteStyles returns the default focused style of the Note field.
func DefaultNoteStyles() (NoteStyle, NoteStyle) {
	focused := NoteStyle{
		Title: lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Margin(1),
		Body:  lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Margin(1),
	}
	blurred := NoteStyle{
		Title: lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Margin(1),
		Body:  lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Margin(1),
	}
	return focused, blurred
}
