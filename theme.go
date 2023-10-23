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

func (t Theme) copy() Theme {
	return Theme{
		Form:           t.Form.Copy(),
		Group:          t.Group.Copy(),
		FieldSeparator: t.FieldSeparator.Copy(),
		Blurred:        t.Blurred.copy(),
		Focused:        t.Focused.copy(),
	}
}

// FieldStyles are the styles for input fields
type FieldStyles struct {
	Base        lipgloss.Style
	Title       lipgloss.Style
	Description lipgloss.Style
	Help        lipgloss.Style
	Error       lipgloss.Style

	// Select styles.
	SelectSelector lipgloss.Style // Selection indicator
	Option         lipgloss.Style // Select options

	// Multi-select styles.
	MultiSelectSelector lipgloss.Style
	SelectedOption      lipgloss.Style
	SelectedPrefix      lipgloss.Style
	UnselectedOption    lipgloss.Style
	UnselectedPrefix    lipgloss.Style

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

func (f FieldStyles) copy() FieldStyles {
	return FieldStyles{
		Base:                f.Base.Copy(),
		Title:               f.Title.Copy(),
		Description:         f.Description.Copy(),
		Help:                f.Help.Copy(),
		Error:               f.Error.Copy(),
		SelectSelector:      f.SelectSelector.Copy(),
		Option:              f.Option.Copy(),
		MultiSelectSelector: f.MultiSelectSelector.Copy(),
		SelectedOption:      f.SelectedOption.Copy(),
		SelectedPrefix:      f.SelectedPrefix.Copy(),
		UnselectedOption:    f.UnselectedOption.Copy(),
		UnselectedPrefix:    f.UnselectedPrefix.Copy(),
		Cursor:              f.Cursor.Copy(),
		Placeholder:         f.Placeholder.Copy(),
		FocusedButton:       f.FocusedButton.Copy(),
		BlurredButton:       f.BlurredButton.Copy(),
		Card:                f.Card.Copy(),
		Next:                f.Next.Copy(),
	}
}

// NewBaseTheme returns a new base theme with general styles to be inherited by
// other themes.
func NewBaseTheme() *Theme {
	var t Theme

	button := lipgloss.NewStyle().Padding(0, 1).Margin(0, 1)

	t.Focused = FieldStyles{
		Base: lipgloss.NewStyle().
			PaddingLeft(1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true),
		SelectSelector:      lipgloss.NewStyle().SetString("> "),
		MultiSelectSelector: lipgloss.NewStyle().SetString("> "),
		SelectedPrefix:      lipgloss.NewStyle().SetString("[•] "),
		UnselectedPrefix:    lipgloss.NewStyle().SetString("[ ] "),
		FocusedButton: button.Copy().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("7")),
		BlurredButton: button.Copy().
			Foreground(lipgloss.Color("7")).
			Background(lipgloss.Color("0")),
	}

	t.Blurred = t.Focused.copy()
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.MultiSelectSelector = lipgloss.NewStyle().SetString("  ")

	return &t
}
