package huh

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Regression for #719: auto-sizing to the tallest group left a large gap
// between a short group's fields and its help footer.
func TestWindowSizeKeepsShortGroupHelpAdjacent(t *testing.T) {
	tall := NewGroup(
		NewNote().
			Title("Welcome").
			Description(strings.Repeat("This is a long welcome note.\n", 12)),
	)
	short := NewGroup(
		NewInput().Title("Name").Description("Your name"),
	)

	f := NewForm(tall, short)
	f.Update(f.Init())
	f.Update(tea.WindowSizeMsg{Width: 80, Height: 60})

	tallH := tall.rawHeight()
	shortNatural := lipgloss.Height(short.Content()) + short.titleFooterHeight()

	if short.height >= tallH && tallH > shortNatural+3 {
		t.Fatalf("short group height=%d was stretched to tall group size (tall rawHeight=%d, short natural≈%d)",
			short.height, tallH, shortNatural)
	}

	// Navigate to the short group and ensure help sits near the input.
	f = batchUpdate(f, f.NextGroup()).(*Form)
	view := viewModel(f)

	helpIdx := strings.Index(view, "enter submit")
	if helpIdx < 0 {
		helpIdx = strings.Index(view, "submit")
	}
	if helpIdx < 0 {
		t.Fatalf("help not found in view:\n%s", view)
	}

	nameIdx := strings.Index(view, "Name")
	if nameIdx < 0 {
		t.Fatalf("input title not found in view:\n%s", view)
	}
	if nameIdx > helpIdx {
		t.Fatalf("expected input before help:\n%s", view)
	}

	between := view[nameIdx:helpIdx]
	blankRuns := 0
	maxBlank := 0
	for _, line := range strings.Split(between, "\n") {
		if strings.TrimSpace(line) == "" {
			blankRuns++
			if blankRuns > maxBlank {
				maxBlank = blankRuns
			}
		} else {
			blankRuns = 0
		}
	}
	if maxBlank > 2 {
		t.Fatalf("expected at most 2 consecutive blank lines between content and help, got %d:\n%s", maxBlank, view)
	}
}
