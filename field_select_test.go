package huh

import "testing"

func TestSelectCachedOptionsUpdateValue(t *testing.T) {
	for _, tt := range []struct {
		name     string
		options  []Option[string]
		selected int
		want     string
	}{
		{"first", NewOptions("Boston", "New York"), 0, "Boston"},
		{"nonzero", NewOptions("Boston", "New York"), 1, "New York"},
		{"clamped", NewOptions("Boston", "New York"), 2, "New York"},
		{"preselected", []Option[string]{NewOption("Boston", "Boston").Selected(true), NewOption("New York", "New York")}, 1, "New York"},
		{"empty", []Option[string]{}, 1, "Toronto"},
		{"nil", nil, 1, "Toronto"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			country := "US"
			city := "Toronto"
			s := NewSelect[string]().Value(&city).OptionsFunc(func() []Option[string] {
				t.Fatal("cached options should not be recomputed")
				return nil
			}, &country)
			s.options.cache[hash(&country)] = tt.options
			s.selected = tt.selected

			s.Update(updateFieldMsg{})

			if city != tt.want {
				t.Errorf("bound value = %q, want %q", city, tt.want)
			}
			if got := s.GetValue(); got != tt.want {
				t.Errorf("GetValue() = %v, want %q", got, tt.want)
			}
			hovered, ok := s.Hovered()
			if len(tt.options) == 0 {
				if ok {
					t.Errorf("Hovered() = %q, want no option", hovered)
				}
			} else if !ok || hovered != tt.want {
				t.Errorf("Hovered() = %q, %v, want %q, true", hovered, ok, tt.want)
			}
		})
	}
}
