package huh

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

func runAccessibleForm(input string, groups ...*Group) error {
	return NewForm(groups...).
		WithAccessible(true).
		WithInput(strings.NewReader(input)).
		WithOutput(&bytes.Buffer{}).
		Run()
}

func TestAccessibleFormInputEndingEarlyAborts(t *testing.T) {
	var name string
	err := runAccessibleForm("", NewGroup(
		NewInput().Title("Name").Value(&name).Validate(func(s string) error {
			if s == "" {
				return errors.New("required")
			}
			return nil
		}),
	))

	if !errors.Is(err, ErrUserAborted) {
		t.Fatalf("expected ErrUserAborted, got %v", err)
	}
	if !errors.Is(err, io.EOF) {
		t.Fatalf("expected the abort to wrap io.EOF, got %v", err)
	}
	if name != "" {
		t.Fatalf("expected no value to be set, got %q", name)
	}
}

func TestAccessibleFormStopsAtTheFieldThatRanOut(t *testing.T) {
	var first, second string
	err := runAccessibleForm("one\n", NewGroup(
		NewInput().Title("First").Value(&first),
		NewInput().Title("Second").Value(&second),
	))

	if !errors.Is(err, ErrUserAborted) {
		t.Fatalf("expected ErrUserAborted, got %v", err)
	}
	if first != "one" {
		t.Fatalf("expected the answered field to keep its value, got %q", first)
	}
}

func TestAccessibleSelectWithClosedInputAborts(t *testing.T) {
	var choice string
	err := runAccessibleForm("", NewGroup(
		NewSelect[string]().Options(NewOptions("a", "b")...).Value(&choice),
	))

	if !errors.Is(err, ErrUserAborted) {
		t.Fatalf("expected ErrUserAborted, got %v", err)
	}
}

func TestAccessiblePromptsEachReadTheirOwnAnswer(t *testing.T) {
	var (
		name  string
		agree bool
		pick  string
	)
	err := runAccessibleForm("ada\ny\n2\n", NewGroup(
		NewInput().Title("Name").Value(&name),
		NewConfirm().Title("Agree").Value(&agree),
		NewSelect[string]().Options(NewOptions("a", "b")...).Value(&pick),
	))
	if err != nil {
		t.Fatal(err)
	}
	if name != "ada" || !agree || pick != "b" {
		t.Fatalf("expected ada/true/b, got %q/%v/%q", name, agree, pick)
	}
}

func TestAccessibleEnterAcceptsADefaultTheValidatorAllows(t *testing.T) {
	name := "ada"
	err := runAccessibleForm("\n", NewGroup(
		NewInput().Title("Name").Value(&name).Validate(func(s string) error {
			if s == "" {
				return errors.New("required")
			}
			return nil
		}),
	))
	if err != nil {
		t.Fatal(err)
	}
	if name != "ada" {
		t.Fatalf("expected the default to stand, got %q", name)
	}
}
