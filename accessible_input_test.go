package huh

import (
	"bytes"
	"strings"
	"testing"
)

func TestAccessiblePromptsEachReadTheirOwnAnswer(t *testing.T) {
	var (
		name  string
		agree bool
		pick  string
	)
	err := NewForm(NewGroup(
		NewInput().Title("Name").Value(&name),
		NewConfirm().Title("Agree").Value(&agree),
		NewSelect[string]().Options(NewOptions("a", "b")...).Value(&pick),
	)).
		WithAccessible(true).
		WithInput(strings.NewReader("ada\ny\n2\n")).
		WithOutput(&bytes.Buffer{}).
		Run()
	if err != nil {
		t.Fatal(err)
	}
	if name != "ada" || !agree || pick != "b" {
		t.Fatalf("expected ada/true/b, got %q/%v/%q", name, agree, pick)
	}
}
