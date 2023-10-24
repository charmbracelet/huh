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
	Base         lipgloss.Style
	Title        lipgloss.Style
	Description  lipgloss.Style
	Help         lipgloss.Style
	Error        lipgloss.Style
	ErrorMessage lipgloss.Style

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
	TextInput TextInputStyles

	// Confirm styles.
	FocusedButton lipgloss.Style
	BlurredButton lipgloss.Style

	// Card styles.
	Card lipgloss.Style
	Next lipgloss.Style
}

type TextInputStyles struct {
	Cursor      lipgloss.Style
	Placeholder lipgloss.Style
	Prompt      lipgloss.Style
	Text        lipgloss.Style
}

func (t TextInputStyles) copy() TextInputStyles {
	return TextInputStyles{
		Cursor:      t.Cursor.Copy(),
		Placeholder: t.Placeholder.Copy(),
		Prompt:      t.Prompt.Copy(),
		Text:        t.Text.Copy(),
	}
}

func (f FieldStyles) copy() FieldStyles {
	return FieldStyles{
		Base:                f.Base.Copy(),
		Title:               f.Title.Copy(),
		Description:         f.Description.Copy(),
		Help:                f.Help.Copy(),
		Error:               f.Error.Copy(),
		ErrorMessage:        f.Error.Copy(),
		SelectSelector:      f.SelectSelector.Copy(),
		Option:              f.Option.Copy(),
		MultiSelectSelector: f.MultiSelectSelector.Copy(),
		SelectedOption:      f.SelectedOption.Copy(),
		SelectedPrefix:      f.SelectedPrefix.Copy(),
		UnselectedOption:    f.UnselectedOption.Copy(),
		UnselectedPrefix:    f.UnselectedPrefix.Copy(),
		FocusedButton:       f.FocusedButton.Copy(),
		BlurredButton:       f.BlurredButton.Copy(),
		TextInput:           f.TextInput.copy(),
		Card:                f.Card.Copy(),
		Next:                f.Next.Copy(),
	}
}

// NewBaseTheme returns a new base theme with general styles to be inherited by
// other themes.
func NewBaseTheme() *Theme {
	var t Theme

	button := lipgloss.NewStyle().Padding(0, 2).MarginRight(1)

	t.Focused = FieldStyles{
		Base: lipgloss.NewStyle().
			PaddingLeft(1).
			BorderStyle(lipgloss.ThickBorder()).
			BorderLeft(true),
		ErrorMessage: lipgloss.NewStyle().
			SetString("* "),
		SelectSelector: lipgloss.NewStyle().
			SetString("> "),
		MultiSelectSelector: lipgloss.NewStyle().
			SetString("> "),
		SelectedPrefix: lipgloss.NewStyle().
			SetString("[•] "),
		UnselectedPrefix: lipgloss.NewStyle().
			SetString("[ ] "),
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

// NewCharmTheme returns a new theme based on the Charm color scheme.
func NewCharmTheme() *Theme {
	t := NewBaseTheme().copy()

	f := &t.Focused
	f.Base = f.Base.BorderForeground(lipgloss.Color("8"))
	f.Title.Foreground(lipgloss.Color("99")).Bold(true)
	f.Description.Foreground(lipgloss.Color("240"))
	f.Help.Foreground(lipgloss.Color("8"))
	f.Error.Foreground(lipgloss.Color("9"))
	f.ErrorMessage.Foreground(lipgloss.Color("9"))
	f.SelectSelector.Foreground(lipgloss.Color("212"))
	f.Option.Foreground(lipgloss.Color("7"))
	f.MultiSelectSelector.Foreground(lipgloss.Color("212"))
	f.SelectedOption.Foreground(lipgloss.Color("212"))
	f.SelectedPrefix.Foreground(lipgloss.Color("212"))
	f.UnselectedOption.Foreground(lipgloss.Color("7"))
	f.FocusedButton.Foreground(lipgloss.Color("#ffffd7")).Background(lipgloss.Color("212"))
	f.BlurredButton.Foreground(lipgloss.Color("7")).Background(lipgloss.Color("0"))

	f.TextInput.Cursor.Foreground(lipgloss.Color("212"))
	f.TextInput.Placeholder.Foreground(lipgloss.Color("8"))
	f.TextInput.Prompt.Foreground(lipgloss.Color("212"))

	t.Blurred = f.copy()
	t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())

	return &t
}

// NewDraculaTheme returns a new theme based on the Dracula color scheme.
func NewDraculaTheme() *Theme {
	t := NewBaseTheme().copy()

	f := &t.Focused
	f.Base.BorderForeground(lipgloss.Color("8"))
	f.Title.Foreground(lipgloss.Color("#bd93f9"))
	f.Description.Foreground(lipgloss.Color("#bd93f9"))
	f.Help.Foreground(lipgloss.Color("8"))
	f.Error.Foreground(lipgloss.Color("#ff5555"))
	f.ErrorMessage.Foreground(lipgloss.Color("#ff5555"))
	f.SelectSelector.Foreground(lipgloss.Color("#f1fa8c"))
	f.Option.Foreground(lipgloss.Color("#f8f8f2"))
	f.MultiSelectSelector.Foreground(lipgloss.Color("#f1fa8c"))
	f.SelectedOption.Foreground(lipgloss.Color("#f1fa8c"))
	f.SelectedPrefix.Foreground(lipgloss.Color("#f1fa8c"))
	f.UnselectedOption.Foreground(lipgloss.Color("#f8f8f2"))
	f.FocusedButton.Foreground(lipgloss.Color("0")).Background(lipgloss.Color("#f1fa8c")).Bold(true)
	f.BlurredButton.Foreground(lipgloss.Color("7")).Background(lipgloss.Color("0"))

	f.TextInput.Cursor.Foreground(lipgloss.Color("#f1fa8c"))
	f.TextInput.Placeholder.Foreground(lipgloss.Color("8"))
	f.TextInput.Prompt.Foreground(lipgloss.Color("#f1fa8c"))

	t.Blurred = f.copy()
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())

	return &t
}

// NewBase16Theme returns a new theme based on the base16 color scheme.
func NewBase16Theme() *Theme {
	t := NewBaseTheme().copy()

	f := &t.Focused
	f.Base.BorderForeground(lipgloss.Color("8"))
	f.Title.Foreground(lipgloss.Color("6"))
	f.Description.Foreground(lipgloss.Color("8"))
	f.Help.Foreground(lipgloss.Color("8"))
	f.Error.Foreground(lipgloss.Color("9"))
	f.ErrorMessage.Foreground(lipgloss.Color("9"))
	f.SelectSelector.Foreground(lipgloss.Color("6"))
	f.Option.Foreground(lipgloss.Color("7"))
	f.MultiSelectSelector.Foreground(lipgloss.Color("6"))
	f.SelectedOption.Foreground(lipgloss.Color("6"))
	f.SelectedPrefix.Foreground(lipgloss.Color("6"))
	f.UnselectedOption.Foreground(lipgloss.Color("7"))
	f.FocusedButton.Foreground(lipgloss.Color("0")).Background(lipgloss.Color("6"))
	f.BlurredButton.Foreground(lipgloss.Color("7")).Background(lipgloss.Color("0"))

	f.TextInput.Cursor.Foreground(lipgloss.Color("6"))
	f.TextInput.Placeholder.Foreground(lipgloss.Color("8"))
	f.TextInput.Prompt.Foreground(lipgloss.Color("6"))

	t.Blurred = f.copy()
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())

	return &t
}
