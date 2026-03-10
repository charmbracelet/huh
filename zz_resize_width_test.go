package huh

import "testing"

func TestSelectWithWidthUpdatesViewportWidth(t *testing.T) {
	f := NewSelect[string]().
		Title("Pick one").
		Options(
			NewOption("Option 1", "1"),
			NewOption("Option 2", "2"),
		)

	f.WithWidth(18)
	if got, want := f.viewport.Width(), 18; got != want {
		t.Fatalf("viewport width after first WithWidth = %d, want %d", got, want)
	}

	f.WithWidth(42)
	if got, want := f.viewport.Width(), 42; got != want {
		t.Fatalf("viewport width after resize WithWidth = %d, want %d", got, want)
	}
}

func TestMultiSelectWithWidthUpdatesViewportWidth(t *testing.T) {
	f := NewMultiSelect[string]().
		Title("Pick many").
		Options(
			NewOption("Option 1", "1"),
			NewOption("Option 2", "2"),
		)

	f.WithWidth(20)
	if got, want := f.viewport.Width(), 20; got != want {
		t.Fatalf("viewport width after first WithWidth = %d, want %d", got, want)
	}

	f.WithWidth(44)
	if got, want := f.viewport.Width(), 44; got != want {
		t.Fatalf("viewport width after resize WithWidth = %d, want %d", got, want)
	}
}
