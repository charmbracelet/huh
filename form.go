package huh

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
)

// Form is a collection of groups that are displayed one at a time on a "page".
//
// The form can navigate between groups and is complete once all the groups are
// complete.
type Form struct {
	// collection of groups
	groups []*Group

	// navigation
	paginator paginator.Model

	// whether or not to use bubble tea rendering for accessibility
	// purposes, if true, the form will render with basic prompting primitives
	// to be more accessible to screen readers.
	accessible bool

	quitting bool

	// options
	theme  *Theme
	keymap *KeyMap
}

// NewForm returns a form with the given groups and default themes and
// keybindings.
//
// Use With* methods to customize the form with options, such as setting
// different themes and keybindings.
func NewForm(groups ...*Group) *Form {
	p := paginator.New()
	p.SetTotalPages(len(groups))

	f := Form{
		groups:    groups,
		paginator: p,
		theme:     NewCharmTheme(),
		keymap:    NewDefaultKeyMap(),
	}

	// NB: If dynamic forms come into play this will need to be applied when
	// groups and fields are added.
	f.applyThemeToChildren()
	f.applyKeymapToChildren()

	return &f
}

// Field is a primitive of a form.
//
// A field represents a single input control on a form such as a text input,
// confirm button, select option, etc...
//
// Each field implements the Bubble Tea Model interface.
type Field interface {
	// Bubble Tea Model
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string

	// Bubble Tea Events
	Blur() tea.Cmd
	Focus() tea.Cmd

	// Errors and Validation
	Error() error

	// Accessible sets whether the field should run in accessible mode.
	Accessible(bool) Field

	// Run runs the field individually.
	Run() error

	// Theme sets the theme on a field.
	Theme(*Theme) Field

	// KeyMap sets the keymap on a field.
	KeyMap(*KeyMap) Field
	KeyBinds() []key.Binding
}

// nextGroupMsg is a message to move to the next group.
type nextGroupMsg struct{}

// prevGroupMsg is a message to move to the previous group.
type prevGroupMsg struct{}

// nextGroup is the command to move to the next group.
func nextGroup() tea.Msg {
	return nextGroupMsg{}
}

// prevGroup is the command to move to the previous group.
func prevGroup() tea.Msg {
	return prevGroupMsg{}
}

// WithAccessible sets the form to run in accessible mode to avoid redrawing the
// views which makes it easier for screen readers to read and describe the form.
//
// This avoids using the Bubble Tea renderer and instead simply uses basic
// terminal prompting to gather input which degrades the user experience but
// provides accessibility.
func (f *Form) WithAccessible(accessible bool) *Form {
	f.accessible = accessible
	return f
}

// WithHelp sets whether or not the form should show help.
//
// This allows the form groups and field to show what keybindings are available
// to the user.
func (f *Form) WithHelp(v bool) *Form {
	for _, group := range f.groups {
		group.ShowHelp(v)
	}
	return f
}

// WithTheme sets the theme on a form.
//
// This allows all groups and fields to be themed consistently, however themes
// can be applied to each group and field individually for more granular
// control.
func (f *Form) WithTheme(theme *Theme) *Form {
	if theme != nil {
		f.theme = theme
		f.applyThemeToChildren()
	}

	return f
}

// applyThemeToChildren applies the form's theme to all children (groups and
// fields).
func (f *Form) applyThemeToChildren() {
	if f.theme == nil {
		return
	}
	for _, group := range f.groups {
		group.Theme(f.theme)
		for _, field := range group.fields {
			field.Theme(f.theme)
		}
	}
}

// KeyMap sets the keymap on a form.
func (f *Form) WithKeyMap(keymap *KeyMap) *Form {
	if keymap == nil {
		return f
	}

	f.keymap = keymap
	f.applyKeymapToChildren()
	return f
}

// applyKeymapToChildren applies the form's keymap to all children (groups and
// fields).
func (f *Form) applyKeymapToChildren() {
	if f.keymap == nil {
		return
	}
	for _, group := range f.groups {
		group.KeyMap(f.keymap)
		for _, field := range group.fields {
			field.KeyMap(f.keymap)
		}
	}
}

// Init initializes the form.
func (f *Form) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, group := range f.groups {
		cmds = append(cmds, group.Init())
	}
	return tea.Batch(cmds...)
}

// Update updates the form.
func (f *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	page := f.paginator.Page
	group := f.groups[page]

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, f.keymap.Quit):
			f.quitting = true
			return f, tea.Quit
		}

	case nextGroupMsg:
		if len(group.Errors()) > 0 {
			return f, nil
		}

		if f.paginator.OnLastPage() {
			f.quitting = true
			return f, tea.Quit
		}
		f.paginator.NextPage()

	case prevGroupMsg:
		if len(group.Errors()) > 0 {
			return f, nil
		}
		f.paginator.PrevPage()
	}

	m, cmd := group.Update(msg)
	f.groups[page] = m.(*Group)

	return f, cmd
}

// View renders the form.
func (f *Form) View() string {
	if f.quitting {
		return ""
	}

	return f.groups[f.paginator.Page].View()
}

// Run runs the form.
func (f *Form) Run() error {
	if len(f.groups) == 0 {
		return nil
	}

	if f.accessible {
		return f.runAccessible()
	}

	return f.run()
}

// run runs the form in normal mode.
func (f *Form) run() error {
	_, err := tea.NewProgram(f).Run()
	return err
}

// runAccessible runs the form in accessible mode.
func (f *Form) runAccessible() error {
	for _, group := range f.groups {
		for _, field := range group.fields {
			field.Init()
			field.Focus()
			field.Accessible(true).Run()
		}
	}

	return nil
}
