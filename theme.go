package huh

import "github.com/charmbracelet/lipgloss"

// Theme is a collection of styles for components of the form.
// Themes can be applied to a form using the WithTheme option.
type Theme struct {
	Form           lipgloss.Style
	Group          lipgloss.Style
	FieldSeparator lipgloss.Style
	Blurred        FieldStyles
	Focused        FieldStyles
}

// copy returns a copy of a theme with all children styles copied.
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
	Base           lipgloss.Style
	Title          lipgloss.Style
	Description    lipgloss.Style
	Help           lipgloss.Style // TODO: apply help coloring in theme to help bubble
	ErrorIndicator lipgloss.Style
	ErrorMessage   lipgloss.Style

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

// TextInputStyles are the styles for text inputs.
type TextInputStyles struct {
	Cursor      lipgloss.Style
	Placeholder lipgloss.Style
	Prompt      lipgloss.Style
	Text        lipgloss.Style
}

// copy returns a copy of a TextInputStyles with all children styles copied.
func (t TextInputStyles) copy() TextInputStyles {
	return TextInputStyles{
		Cursor:      t.Cursor.Copy(),
		Placeholder: t.Placeholder.Copy(),
		Prompt:      t.Prompt.Copy(),
		Text:        t.Text.Copy(),
	}
}

// copy returns a copy of a FieldStyles with all children styles copied.
func (f FieldStyles) copy() FieldStyles {
	return FieldStyles{
		Base:                f.Base.Copy(),
		Title:               f.Title.Copy(),
		Description:         f.Description.Copy(),
		Help:                f.Help.Copy(),
		ErrorIndicator:      f.ErrorIndicator.Copy(),
		ErrorMessage:        f.ErrorMessage.Copy(),
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

	t.FieldSeparator = lipgloss.NewStyle().SetString("\n\n")
	button := lipgloss.NewStyle().Padding(0, 2).MarginRight(1)

	// Focused styles.
	f := &t.Focused
	f.Base = lipgloss.NewStyle().
		PaddingLeft(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderLeft(true)
	f.ErrorIndicator = lipgloss.NewStyle().
		SetString(" *")
	f.ErrorMessage = lipgloss.NewStyle().
		SetString(" *")
	f.SelectSelector = lipgloss.NewStyle().
		SetString("> ")
	f.MultiSelectSelector = lipgloss.NewStyle().
		SetString("> ")
	f.SelectedPrefix = lipgloss.NewStyle().
		SetString("[•] ")
	f.UnselectedPrefix = lipgloss.NewStyle().
		SetString("[ ] ")
	f.FocusedButton = button.Copy().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("7"))
	f.BlurredButton = button.Copy().
		Foreground(lipgloss.Color("7")).
		Background(lipgloss.Color("0"))
	f.TextInput.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	f.Help = lipgloss.NewStyle().
		PaddingLeft(1)

	// Blurred styles.
	t.Blurred = f.copy()
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.MultiSelectSelector = lipgloss.NewStyle().SetString("  ")

	return &t
}

// NewCharmTheme returns a new theme based on the Charm color scheme.
func NewCharmTheme() *Theme {
	t := NewBaseTheme().copy()

	var (
		normalFg = lipgloss.AdaptiveColor{Light: "235", Dark: "252"}
		indigo   = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
		cream    = lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}
		fuchsia  = lipgloss.Color("#F780E2")
		green    = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
		red      = lipgloss.AdaptiveColor{Light: "#FF4672", Dark: "#ED567A"}
	)

	f := &t.Focused
	f.Base = f.Base.BorderForeground(lipgloss.Color("238"))
	f.Title.Foreground(indigo).Bold(true)
	f.Description.Foreground(lipgloss.AdaptiveColor{Light: "", Dark: "243"})
	f.Help.Foreground(lipgloss.Color("8"))
	f.ErrorIndicator.Foreground(red)
	f.ErrorMessage.Foreground(red)
	f.SelectSelector.Foreground(fuchsia)
	f.Option.Foreground(normalFg)
	f.MultiSelectSelector.Foreground(fuchsia)
	f.SelectedOption.Foreground(green)
	f.SelectedPrefix = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#02CF92", Dark: "#02A877"}).SetString("[✓] ")
	f.UnselectedPrefix = lipgloss.NewStyle().SetString("[ ] ")
	f.UnselectedOption.Foreground(normalFg)
	f.FocusedButton.Foreground(cream).Background(fuchsia)
	f.Next = f.FocusedButton.Copy()
	f.BlurredButton.Foreground(normalFg).Background(lipgloss.AdaptiveColor{Light: "252", Dark: "237"})

	f.TextInput.Cursor.Foreground(green)
	f.TextInput.Placeholder.Foreground(lipgloss.AdaptiveColor{Light: "248", Dark: "238"})
	f.TextInput.Prompt.Foreground(fuchsia)

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
	f.ErrorIndicator.Foreground(lipgloss.Color("#ff5555"))
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
	f.ErrorIndicator.Foreground(lipgloss.Color("9"))
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
	t.Blurred.Title.Foreground(lipgloss.Color("8"))
	t.Blurred.TextInput.Prompt.Foreground(lipgloss.Color("8"))
	t.Blurred.TextInput.Text.Foreground(lipgloss.Color("8"))

	return &t
}
