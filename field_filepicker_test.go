package huh

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func errRejectFile(string) error { return errors.New("choose an approved file") }

func newTestFilePicker(t *testing.T, validate func(string) error) *FilePicker {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	fp := NewFilePicker().
		CurrentDirectory(dir).
		Validate(validate)
	fp.WithKeyMap(NewDefaultKeyMap())
	fp.Focus()
	return fp
}

// emitsNextField reports whether cmd produces a nextFieldMsg.
func emitsNextField(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(nextFieldMsg)
	return ok
}

func emitsPrevField(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(prevFieldMsg)
	return ok
}

func TestFilePickerNextBlockedByValidation(t *testing.T) {
	fp := newTestFilePicker(t, errRejectFile)
	_, cmd := fp.Update(codeKeypress(tea.KeyTab))
	if emitsNextField(cmd) {
		t.Error("tab advanced to the next field despite a rejecting Validate")
	}
	if fp.Error() == nil {
		t.Error("expected validation error to be set after rejected forward navigation")
	}
}

func TestFilePickerNextAllowedByValidation(t *testing.T) {
	fp := newTestFilePicker(t, func(string) error { return nil })
	_, cmd := fp.Update(codeKeypress(tea.KeyTab))
	if !emitsNextField(cmd) {
		t.Error("tab did not advance even though validation passed")
	}
	if fp.Error() != nil {
		t.Errorf("unexpected error: %v", fp.Error())
	}
}

func TestFilePickerEscapeBlockedByValidation(t *testing.T) {
	fp := newTestFilePicker(t, errRejectFile)
	fp.Picking(true)
	_, cmd := fp.Update(codeKeypress(tea.KeyEscape))
	if emitsNextField(cmd) {
		t.Error("escape while picking advanced despite a rejecting Validate")
	}
	if fp.Error() == nil {
		t.Error("expected validation error after escape with rejecting Validate")
	}
}

func TestFilePickerSelectionBlockedByValidation(t *testing.T) {
	fp := newTestFilePicker(t, errRejectFile)
	fp.Picking(true)
	// enter selects the highlighted file (the only entry in the temp dir).
	_, cmd := fp.Update(codeKeypress(tea.KeyEnter))
	if got := fp.accessor.Get(); got == "" {
		t.Fatal("expected the file selection to be recorded")
	}
	if emitsNextField(cmd) {
		t.Error("selecting a file advanced despite a rejecting Validate")
	}
	if fp.Error() == nil {
		t.Error("expected validation error after selecting a file with rejecting Validate")
	}
}

func TestFilePickerSelectionAllowedByValidation(t *testing.T) {
	fp := newTestFilePicker(t, func(string) error { return nil })
	fp.Picking(true)
	_, cmd := fp.Update(codeKeypress(tea.KeyEnter))
	if got := fp.accessor.Get(); got == "" {
		t.Fatal("expected the file selection to be recorded")
	}
	if !emitsNextField(cmd) {
		t.Error("selecting a file did not advance even though validation passed")
	}
}

func TestFilePickerPrevSkipsValidation(t *testing.T) {
	fp := newTestFilePicker(t, errRejectFile)
	_, cmd := fp.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyTab, Mod: tea.ModShift}))
	if !emitsPrevField(cmd) {
		t.Error("shift+tab should navigate backward without validation")
	}
}

func TestFilePickerErrorSurvivesUnrelatedMessages(t *testing.T) {
	fp := newTestFilePicker(t, errRejectFile)
	// Blur validates and records the rejection, as when the group moves away.
	fp.Blur()
	if fp.Error() == nil {
		t.Fatal("expected Blur to record the validation error")
	}
	// The group forwards unrelated background and refresh messages to every
	// field; these must not clear a recorded validation error.
	fp.Update(tea.BackgroundColorMsg{})
	if fp.Error() == nil {
		t.Error("validation error cleared by a background color message")
	}
	fp2 := newTestFilePicker(t, errRejectFile)
	fp2.Blur()
	fp2.Update(updateFieldMsg{})
	if fp2.Error() == nil {
		t.Error("validation error cleared by the group's field refresh")
	}
}
