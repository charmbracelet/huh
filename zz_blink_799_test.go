package huh

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

// TestIssue799BlinkAfterView exercises the huh fix for issue #799: calling
// View() on a focused Input must NOT kill the in-flight cursor BlinkMsg
// chain.
//
// The pre-fix View() unconditionally called textinput.SetStyles, which
// internally calls cursor.SetMode(CursorBlink), which bumps the cursor's
// blinkTag. A Blink cmd that was captured with the old blinkTag now
// produces a BlinkMsg whose tag is stale, and the cursor rejects it,
// silently halting the blink chain.
//
// The fix introduces a `stylesDirty` flag on Input so View() only calls
// SetStyles when something actually changed (Focus, Blur, theme switch,
// background-color msg). On a plain View() call with no state change
// SetStyles is skipped, and the existing Blink chain is preserved.
func TestIssue799BlinkAfterView(t *testing.T) {
	input := NewInput()

	// Focus returns the initial Blink cmd.
	cmd := input.Focus()
	if cmd == nil {
		t.Fatal("expected initial Blink cmd from Focus()")
	}

	// Simulate the render loop: a View() without any state change. The
	// huh-level fix must NOT touch the underlying textinput (no
	// SetStyles, no cursor.SetMode, no blinkTag bump).
	for i := 0; i < 5; i++ {
		_ = input.View() // should be a no-op for blinkTag purposes

		msg := cmd()
		if msg == nil {
			t.Fatalf("iter %d: Blink cmd returned nil", i)
		}
		_, next := input.Update(msg)
		if next == nil {
			t.Fatalf("iter %d: Blink chain died after View() (issue #799): msg=%T", i, msg)
		}
		cmd = next
	}
}

// silence unused-import linter for tea when the test above doesn't reference
// tea directly. The huh types use tea.Msg under the hood.
var _ tea.Cmd
