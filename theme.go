package huh

import "github.com/charmbracelet/lipgloss"

// Theme is the style for a form.
type Theme struct {
	Form           lipgloss.Style
	Group          lipgloss.Style
	FieldSeparator lipgloss.Style
	Blurred        FieldStyles
	Focused        FieldStyles
}

// FieldStyles are the styles for input fields
type FieldStyles struct {
	Base        lipgloss.Style
	Title       lipgloss.Style
	Description lipgloss.Style
	Help        lipgloss.Style
	Error       lipgloss.Style

	// Select and multi-select styles.
	Selector lipgloss.Style // Selection indicator in selects and multi-selects
	Option   lipgloss.Style // Select options

	// Multi-select styles.
	SelectedOption   lipgloss.Style
	SelectedPrefix   lipgloss.Style
	UnselectedOption lipgloss.Style
	UnselectedPrefix lipgloss.Style

	// Textinput and teatarea styles.
	Cursor      lipgloss.Style
	Placeholder lipgloss.Style

	// Confirm styles.
	FocusedButton lipgloss.Style
	BlurredButton lipgloss.Style

	// Card styles.
	Card lipgloss.Style
	Next lipgloss.Style
}

// NewBaseTheme returns a new base theme with general styles to be inherited by
// other themes.
func NewBaseTheme() *Theme {
	var t Theme

	button := lipgloss.NewStyle().Padding(0, 1).Margin(0, 1)

	t.Blurred = FieldStyles{
		Base: lipgloss.NewStyle().
			PaddingLeft(1).
			BorderStyle(lipgloss.HiddenBorder()).
			BorderLeft(true),
		Selector:         lipgloss.NewStyle().SetString("> "),
		SelectedPrefix:   lipgloss.NewStyle().SetString("[•] "),
		UnselectedPrefix: lipgloss.NewStyle().SetString("[ ] "),
		FocusedButton: button.Copy().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("7")),
		BlurredButton: button.Copy().
			Foreground(lipgloss.Color("7")).
			Background(lipgloss.Color("0")),
	}

	t.Focused = FieldStyles{
		Base:             t.Blurred.Base.Copy().BorderStyle(lipgloss.NormalBorder()),
		Selector:         t.Blurred.Selector.Copy(),
		SelectedPrefix:   t.Blurred.SelectedPrefix.Copy(),
		UnselectedPrefix: t.Blurred.UnselectedPrefix.Copy(),
		FocusedButton:    t.Blurred.FocusedButton.Copy(),
		BlurredButton:    t.Blurred.BlurredButton.Copy(),
	}

	return &t
}
