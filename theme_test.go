package huh

import "testing"

func TestThemeBaseErrorMessageAlignment(t *testing.T) {
	styles := ThemeBase(false)

	if got := styles.Focused.ErrorIndicator.String(); got != " *" {
		t.Fatalf("ErrorIndicator = %q, want %q", got, " *")
	}

	if got := styles.Focused.ErrorMessage.String(); got != "*" {
		t.Fatalf("ErrorMessage = %q, want %q (no leading space so footer aligns with help)", got, "*")
	}
}
