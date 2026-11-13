package huh

import (
	"image/color"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestInputBlurredTextUsesTheme(t *testing.T) {
	const (
		textColour        = "#11AA22"
		promptColour      = "#334455"
		placeholderColour = "#667788"
	)
	theme := ThemeCharm(true)
	theme.Blurred.TextInput.Text = lipgloss.NewStyle().Foreground(lipgloss.Color(textColour))
	theme.Blurred.TextInput.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color(promptColour))
	theme.Blurred.TextInput.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color(placeholderColour))

	value := "meshns"
	field := NewInput().Title("Namespace").Value(&value)
	field.WithTheme(ThemeFunc(func(bool) *Styles { return theme }))
	field.WithWidth(40)
	field.Blur()
	view := field.View()

	if !strings.Contains(view, value) {
		t.Fatalf("blurred input view %q does not contain the value", view)
	}
	st := field.textinput.Styles()
	assertColour(t, "blurred text", st.Blurred.Text.GetForeground(), lipgloss.Color(textColour))
	assertColour(t, "blurred prompt", st.Blurred.Prompt.GetForeground(), lipgloss.Color(promptColour))
	assertColour(t, "blurred placeholder", st.Blurred.Placeholder.GetForeground(), lipgloss.Color(placeholderColour))
	if coloursEqual(st.Focused.Text.GetForeground(), lipgloss.Color(textColour)) {
		t.Fatalf("blurred render wrote the blurred text colour onto the focused style")
	}
}

func TestInputFocusedTextUsesTheme(t *testing.T) {
	const textColour = "#ABCDEF"
	theme := ThemeCharm(true)
	theme.Focused.TextInput.Text = lipgloss.NewStyle().Foreground(lipgloss.Color(textColour))

	value := "meshns"
	field := NewInput().Value(&value)
	field.WithTheme(ThemeFunc(func(bool) *Styles { return theme }))
	field.WithWidth(40)
	field.Focus()
	field.View()

	assertColour(t, "focused text", field.textinput.Styles().Focused.Text.GetForeground(), lipgloss.Color(textColour))
}

func TestTextCursorUsesForeground(t *testing.T) {
	const (
		foreground = "#11AA22"
		background = "#CC0000"
	)
	theme := ThemeCharm(true)
	theme.Focused.TextInput.Cursor = lipgloss.NewStyle().
		Foreground(lipgloss.Color(foreground)).
		Background(lipgloss.Color(background))

	value := "hello"
	field := NewText().Value(&value)
	field.WithTheme(ThemeFunc(func(bool) *Styles { return theme }))
	field.WithWidth(40)
	field.Focus()
	field.View()

	assertColour(t, "cursor", field.textarea.Styles().Cursor.Color, lipgloss.Color(foreground))
}

func assertColour(t *testing.T, name string, got, want color.Color) {
	t.Helper()
	if !coloursEqual(got, want) {
		t.Fatalf("%s colour = %v, want %v", name, got, want)
	}
}

func coloursEqual(a, b color.Color) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return ar == br && ag == bg && ab == bb && aa == ba
}
