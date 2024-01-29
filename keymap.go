package huh

import "github.com/charmbracelet/bubbles/key"

// KeyMap is the keybindings to navigate the form.
type KeyMap struct {
	Quit key.Binding

	Confirm     ConfirmKeyMap
	File        FileKeyMap
	Input       InputKeyMap
	MultiSelect MultiSelectKeyMap
	Note        NoteKeyMap
	Select      SelectKeyMap
	Text        TextKeyMap
}

// InputKeyMap is the keybindings for input fields.
type InputKeyMap struct {
	AcceptSuggestion key.Binding
	Next             key.Binding
	Prev             key.Binding
	Submit           key.Binding
}

// TextKeyMap is the keybindings for text fields.
type TextKeyMap struct {
	Next    key.Binding
	Prev    key.Binding
	NewLine key.Binding
	Editor  key.Binding
	Submit  key.Binding
}

// SelectKeyMap is the keybindings for select fields.
type SelectKeyMap struct {
	Next        key.Binding
	Prev        key.Binding
	Up          key.Binding
	Down        key.Binding
	Filter      key.Binding
	SetFilter   key.Binding
	ClearFilter key.Binding
	Submit      key.Binding
}

// MultiSelectKeyMap is the keybindings for multi-select fields.
type MultiSelectKeyMap struct {
	Next        key.Binding
	Prev        key.Binding
	Up          key.Binding
	Down        key.Binding
	Toggle      key.Binding
	Filter      key.Binding
	SetFilter   key.Binding
	ClearFilter key.Binding
	Submit      key.Binding
}

// FilePicker is the keybindings for filepicker fields.
type FileKeyMap struct {
	GoToTop  key.Binding
	GoToLast key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Back     key.Binding
	Open     key.Binding
	Select   key.Binding
	Up       key.Binding
	Down     key.Binding
	Prev     key.Binding
	Next     key.Binding
	Submit   key.Binding
}

// NoteKeyMap is the keybindings for note fields.
type NoteKeyMap struct {
	Next   key.Binding
	Prev   key.Binding
	Submit key.Binding
}

// ConfirmKeyMap is the keybindings for confirm fields.
type ConfirmKeyMap struct {
	Next   key.Binding
	Prev   key.Binding
	Toggle key.Binding
	Submit key.Binding
}

// NewDefaultKeyMap returns a new default keymap.
func NewDefaultKeyMap() *KeyMap {
	return &KeyMap{
		Quit: key.NewBinding(key.WithKeys("ctrl+c")),
		Input: InputKeyMap{
			AcceptSuggestion: key.NewBinding(key.WithKeys("ctrl+e"), key.WithHelp("ctrl+e", "complete")),
			Prev:             key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "back")),
			Next:             key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter", "next")),
			Submit:           key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
		},
		File: FileKeyMap{
			GoToTop:  key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "first")),
			GoToLast: key.NewBinding(key.WithKeys("G"), key.WithHelp("G", "last")),
			PageUp:   key.NewBinding(key.WithKeys("K", "pgup"), key.WithHelp("pgup", "page up"), key.WithDisabled()),
			PageDown: key.NewBinding(key.WithKeys("J", "pgdown"), key.WithHelp("pgdown", "page down"), key.WithDisabled()),
			Back:     key.NewBinding(key.WithKeys("h", "backspace", "left", "esc"), key.WithHelp("h", "back")),
			Open:     key.NewBinding(key.WithKeys("l", "right", "enter"), key.WithHelp("enter", "open")),
			Select:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
			Up:       key.NewBinding(key.WithKeys("up", "k", "ctrl+k", "ctrl+p"), key.WithHelp("↑", "up")),
			Down:     key.NewBinding(key.WithKeys("down", "j", "ctrl+j", "ctrl+n"), key.WithHelp("↓", "down")),
			Prev:     key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "back")),
			Next:     key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next")),
			Submit:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
		},
		Text: TextKeyMap{
			Prev:    key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "back")),
			Next:    key.NewBinding(key.WithKeys("tab", "enter"), key.WithHelp("enter", "next")),
			Submit:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
			NewLine: key.NewBinding(key.WithKeys("alt+enter", "ctrl+j"), key.WithHelp("alt+enter / ctrl+j", "new line")),
			Editor:  key.NewBinding(key.WithKeys("ctrl+e"), key.WithHelp("ctrl+e", "open editor")),
		},
		Select: SelectKeyMap{
			Prev:        key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "back")),
			Next:        key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter", "select")),
			Submit:      key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
			Up:          key.NewBinding(key.WithKeys("up", "k", "ctrl+k", "ctrl+p"), key.WithHelp("↑", "up")),
			Down:        key.NewBinding(key.WithKeys("down", "j", "ctrl+j", "ctrl+n"), key.WithHelp("↓", "down")),
			Filter:      key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
			SetFilter:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "set filter"), key.WithDisabled()),
			ClearFilter: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "clear filter"), key.WithDisabled()),
		},
		MultiSelect: MultiSelectKeyMap{
			Prev:        key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "back")),
			Next:        key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter", "confirm")),
			Submit:      key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
			Toggle:      key.NewBinding(key.WithKeys(" ", "x"), key.WithHelp("x", "toggle")),
			Up:          key.NewBinding(key.WithKeys("up", "k", "ctrl+p"), key.WithHelp("↑", "up")),
			Down:        key.NewBinding(key.WithKeys("down", "j", "ctrl+n"), key.WithHelp("↓", "down")),
			Filter:      key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
			SetFilter:   key.NewBinding(key.WithKeys("enter", "esc"), key.WithHelp("esc", "set filter"), key.WithDisabled()),
			ClearFilter: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "clear filter"), key.WithDisabled()),
		},
		Note: NoteKeyMap{
			Prev:   key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "back")),
			Next:   key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter", "next")),
			Submit: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
		},
		Confirm: ConfirmKeyMap{
			Prev:   key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "back")),
			Next:   key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter", "next")),
			Submit: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
			Toggle: key.NewBinding(key.WithKeys("h", "l", "right", "left"), key.WithHelp("←/→", "toggle")),
		},
	}
}
