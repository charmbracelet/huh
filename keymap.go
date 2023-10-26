package huh

import "github.com/charmbracelet/bubbles/key"

// KeyMap is the keybindings to navigate the form.
type KeyMap struct {
	Quit key.Binding

	Input       InputKeyMap
	Text        TextKeyMap
	Select      SelectKeyMap
	MultiSelect MultiSelectKeyMap
	Note        NoteKeyMap
	Confirm     ConfirmKeyMap
}

// InputKeyMap is the keybindings for input fields.
type InputKeyMap struct {
	Next key.Binding
	Prev key.Binding
}

// TextKeyMap is the keybindings for text fields.
type TextKeyMap struct {
	Next    key.Binding
	Prev    key.Binding
	NewLine key.Binding
	Editor  key.Binding
}

// SelectKeyMap is the keybindings for select fields.
type SelectKeyMap struct {
	Next      key.Binding
	Prev      key.Binding
	Up        key.Binding
	Down      key.Binding
	Filter    key.Binding
	SetFilter key.Binding
}

// MultiSelectKeyMap is the keybindings for multi-select fields.
type MultiSelectKeyMap struct {
	Next   key.Binding
	Prev   key.Binding
	Up     key.Binding
	Down   key.Binding
	Toggle key.Binding
}

// NoteKeyMap is the keybindings for note fields.
type NoteKeyMap struct {
	Next key.Binding
	Prev key.Binding
}

// ConfirmKeyMap is the keybindings for confirm fields.
type ConfirmKeyMap struct {
	Next   key.Binding
	Prev   key.Binding
	Toggle key.Binding
}

// NewDefaultKeyMap returns a new default keymap.
func NewDefaultKeyMap() *KeyMap {
	return &KeyMap{
		Quit: key.NewBinding(key.WithKeys("ctrl+c")),
		Input: InputKeyMap{
			Next: key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter", "next")),
			Prev: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "back")),
		},
		Text: TextKeyMap{
			Next:    key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next")),
			Prev:    key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "back")),
			NewLine: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "new line")),
			Editor:  key.NewBinding(key.WithKeys("ctrl+e"), key.WithHelp("ctrl+e", "open editor")),
		},
		Select: SelectKeyMap{
			Next:      key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter", "select")),
			Prev:      key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "back")),
			Up:        key.NewBinding(key.WithKeys("up", "k", "ctrl+p"), key.WithHelp("↑", "up")),
			Down:      key.NewBinding(key.WithKeys("down", "j", "ctrl+n"), key.WithHelp("↓", "down")),
			Filter:    key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
			SetFilter: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "set filter"), key.WithDisabled()),
		},
		MultiSelect: MultiSelectKeyMap{
			Next:   key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter", "confirm")),
			Prev:   key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "back")),
			Toggle: key.NewBinding(key.WithKeys(" ", "x"), key.WithHelp("x", "toggle")),
			Up:     key.NewBinding(key.WithKeys("up", "k", "ctrl+p"), key.WithHelp("↑", "up")),
			Down:   key.NewBinding(key.WithKeys("down", "j", "ctrl+n"), key.WithHelp("↓", "down")),
		},
		Note: NoteKeyMap{
			Next: key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter", "next")),
			Prev: key.NewBinding(key.WithKeys("shift+tab")),
		},
		Confirm: ConfirmKeyMap{
			Next:   key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter", "next")),
			Prev:   key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "back")),
			Toggle: key.NewBinding(key.WithKeys("h", "l", "right", "left"), key.WithHelp("←/→", "toggle")),
		},
	}
}
