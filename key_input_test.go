package huh

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestSanitizeKeyPressForInputClearsNavigationText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		key  tea.Key
	}{
		{name: "tab", key: tea.Key{Code: tea.KeyTab, Text: "\t"}},
		{name: "shift tab", key: tea.Key{Code: tea.KeyTab, Mod: tea.ModShift, Text: "\t"}},
		{name: "enter", key: tea.Key{Code: tea.KeyEnter, Text: "\r"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sanitized, ok := sanitizeKeyPressForInput(tea.KeyPressMsg(tt.key)).(tea.KeyPressMsg)
			if !ok {
				t.Fatalf("expected KeyPressMsg, got %T", sanitized)
			}
			if sanitized.Key().Text != "" {
				t.Fatalf("expected empty text, got %q", sanitized.Key().Text)
			}
		})
	}
}

func TestInputIgnoresWindowsNavigationKeyText(t *testing.T) {
	t.Parallel()

	field := NewInput()
	field.Focus()

	m, _ := field.Update(tea.KeyPressMsg(tea.Key{
		Code: tea.KeyTab,
		Text: "\t",
	}))
	field = m.(*Input)
	if got, _ := field.GetValue().(string); got != "" {
		t.Fatalf("expected empty value after tab, got %q", got)
	}

	m, _ = field.Update(tea.KeyPressMsg(tea.Key{
		Code: tea.KeyEnter,
		Text: "\r",
	}))
	field = m.(*Input)
	if got, _ := field.GetValue().(string); got != "" {
		t.Fatalf("expected empty value after enter, got %q", got)
	}
}

func TestSanitizeKeyPressForInputPassthrough(t *testing.T) {
	t.Parallel()

	paste := tea.PasteMsg{Content: "hello"}
	if got := sanitizeKeyPressForInput(paste); got != paste {
		t.Fatalf("expected paste passthrough, got %#v", got)
	}

	plain := tea.KeyPressMsg(tea.Key{Code: tea.KeyTab})
	if got := sanitizeKeyPressForInput(plain); got != plain {
		t.Fatalf("expected unchanged tab keypress")
	}

	charKey := tea.KeyPressMsg(tea.Key{Code: 'a', Text: "a"})
	if got := sanitizeKeyPressForInput(charKey); got != charKey {
		t.Fatalf("expected unchanged character keypress")
	}
}

func TestSanitizeKeyPressForInputPreservesModifiers(t *testing.T) {
	t.Parallel()

	shiftTab := tea.KeyPressMsg(tea.Key{Code: tea.KeyTab, Mod: tea.ModShift, Text: "\t"})
	sanitized := sanitizeKeyPressForInput(shiftTab).(tea.KeyPressMsg)
	if sanitized.Key().Mod != tea.ModShift {
		t.Fatalf("expected shift modifier to be preserved")
	}
}

func TestTextIgnoresWindowsNavigationKeyText(t *testing.T) {
	t.Parallel()

	field := NewText()
	field.Focus()

	m, _ := field.Update(tea.KeyPressMsg(tea.Key{
		Code: tea.KeyTab,
		Text: "\t",
	}))
	field = m.(*Text)
	if got, _ := field.GetValue().(string); got != "" {
		t.Fatalf("expected empty value after tab, got %q", got)
	}
}
