package huh

import (
	"errors"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
)

// FormState represents the current state of the form.
type FormState int

const (
	// StateNormal is when the user is completing the form.
	StateNormal FormState = iota

	// StateCompleted is when the user has completed the form.
	StateCompleted

	// StateAborted is when the user has aborted the form.
	StateAborted
)

// ErrUserAborted is the error returned when a user exits the form before
// submitting.
var (
	ErrUserAborted = errors.New("user aborted")
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

	// callbacks
	submitCmd tea.Cmd
	cancelCmd tea.Cmd

	State FormState

	// whether or not to use bubble tea rendering for accessibility
	// purposes, if true, the form will render with basic prompting primitives
	// to be more accessible to screen readers.
	accessible bool

	quitting bool
	aborted  bool

	// options
	width  int
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

	f := &Form{
		groups:    groups,
		paginator: p,
		theme:     NewCharmTheme(),
		keymap:    NewDefaultKeyMap(),
		width:     80,
	}

	// NB: If dynamic forms come into play this will need to be applied when
	// groups and fields are added.
	f.WithTheme(f.theme)
	f.WithKeyMap(f.keymap)
	f.WithWidth(f.width)

	return f
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

	// Run runs the field individually.
	Run() error

	// KeyBinds returns help keybindings.
	KeyBinds() []key.Binding

	// WithTheme sets the theme on a field.
	WithTheme(*Theme) Field

	// WithAccessible sets whether the field should run in accessible mode.
	WithAccessible(bool) Field

	// WithKeyMap sets the keymap on a field.
	WithKeyMap(*KeyMap) Field

	// WithWidth sets the width of a field.
	WithWidth(int) Field
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
		group.WithHelp(v)
	}
	return f
}

// WithTheme sets the theme on a form.
//
// This allows all groups and fields to be themed consistently, however themes
// can be applied to each group and field individually for more granular
// control.
func (f *Form) WithTheme(theme *Theme) *Form {
	if theme == nil {
		return f
	}
	f.theme = theme
	for _, group := range f.groups {
		group.WithTheme(theme)
	}
	return f
}

// WithKeyMap sets the keymap on a form.
//
// This allows customization of the form key bindings.
func (f *Form) WithKeyMap(keymap *KeyMap) *Form {
	if keymap == nil {
		return f
	}
	f.keymap = keymap
	for _, group := range f.groups {
		group.WithKeyMap(keymap)
	}
	return f
}

// WithWidth sets the width of a form
//
// This allows all groups and fields to be sized consistently, however width
// can be applied to each group and field individually for more granular
// control.
func (f *Form) WithWidth(width int) *Form {
	if width <= 0 {
		return f
	}
	f.width = width
	for _, group := range f.groups {
		group.WithWidth(width)
	}
	return f
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
			f.aborted = true
			f.quitting = true
			f.State = StateAborted
			return f, f.cancelCmd
		}

	case nextGroupMsg:
		if len(group.Errors()) > 0 {
			return f, nil
		}

		if f.paginator.OnLastPage() {
			f.quitting = true
			f.State = StateCompleted
			return f, f.submitCmd
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
	f.submitCmd = tea.Quit
	f.cancelCmd = tea.Quit

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
	m, err := tea.NewProgram(f).Run()
	if m.(*Form).aborted {
		err = ErrUserAborted
	}
	return err
}

// runAccessible runs the form in accessible mode.
func (f *Form) runAccessible() error {
	for _, group := range f.groups {
		for _, field := range group.fields {
			field.Init()
			field.Focus()
			_ = field.WithAccessible(true).Run()
		}
	}

	return nil
}
