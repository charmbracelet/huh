package huh

import (
	"errors"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
)

const defaultWidth = 80

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

// ErrUserAborted is the error returned when a user exits the form before submitting.
var ErrUserAborted = errors.New("user aborted")

// Form is a collection of groups that are displayed one at a time on a "page".
//
// The form can navigate between groups and is complete once all the groups are
// complete.
type Form struct {
	// collection of groups
	groups []*Group

	results map[string]any

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
	width      int
	height     int
	keymap     *KeyMap
	teaOptions []tea.ProgramOption
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
		keymap:    NewDefaultKeyMap(),
		results:   make(map[string]any),
		teaOptions: []tea.ProgramOption{
			tea.WithOutput(os.Stderr),
		},
	}

	// NB: If dynamic forms come into play this will need to be applied when
	// groups and fields are added.
	f.WithKeyMap(f.keymap)
	f.WithWidth(f.width)
	f.WithHeight(f.height)
	f.UpdateFieldPositions()

	if os.Getenv("TERM") == "dumb" {
		f.WithWidth(defaultWidth)
		f.WithAccessible(true)
	}

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

	// Skip returns whether this input should be skipped or not.
	Skip() bool

	// Zoom returns whether this input should be zoomed or not.
	// Zoom allows the field to take focus of the group / form height.
	Zoom() bool

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

	// WithHeight sets the height of a field.
	WithHeight(int) Field

	// WithPosition tells the field the index of the group and position it is in.
	WithPosition(FieldPosition) Field

	// GetKey returns the field's key.
	GetKey() string

	// GetValue returns the field's value.
	GetValue() any
}

// FieldPosition is positional information about the given field and form.
type FieldPosition struct {
	Group      int
	Field      int
	FirstField int
	LastField  int
	GroupCount int
	FirstGroup int
	LastGroup  int
}

// IsFirst returns whether a field is the form's first field.
func (p FieldPosition) IsFirst() bool {
	return p.Field == p.FirstField && p.Group == p.FirstGroup
}

// IsLast returns whether a field is the form's last field.
func (p FieldPosition) IsLast() bool {
	return p.Field == p.LastField && p.Group == p.LastGroup
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

// WithShowHelp sets whether or not the form should show help.
//
// This allows the form groups and field to show what keybindings are available
// to the user.
func (f *Form) WithShowHelp(v bool) *Form {
	for _, group := range f.groups {
		group.WithShowHelp(v)
	}
	return f
}

// WithShowErrors sets whether or not the form should show errors.
//
// This allows the form groups and fields to show errors when the Validate
// function returns an error.
func (f *Form) WithShowErrors(v bool) *Form {
	for _, group := range f.groups {
		group.WithShowErrors(v)
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
	f.UpdateFieldPositions()
	return f
}

// WithWidth sets the width of a form.
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

// WithHeight sets the height of a form.
func (f *Form) WithHeight(height int) *Form {
	if height <= 0 {
		return f
	}
	f.height = height
	for _, group := range f.groups {
		group.WithHeight(height)
	}
	return f
}

// WithOutput sets the io.Writer to output the form.
func (f *Form) WithOutput(w io.Writer) *Form {
	f.teaOptions = append(f.teaOptions, tea.WithOutput(w))
	return f
}

// WithProgramOptions sets the tea options of the form.
func (f *Form) WithProgramOptions(opts ...tea.ProgramOption) *Form {
	f.teaOptions = opts
	return f
}

// UpdateFieldPositions sets the position on all the fields.
func (f *Form) UpdateFieldPositions() *Form {
	firstGroup := 0
	lastGroup := len(f.groups) - 1

	// determine the first non-hidden group.
	for g := range f.groups {
		if !f.isGroupHidden(g) {
			break
		}
		firstGroup++
	}

	// determine the last non-hidden group.
	for g := len(f.groups) - 1; g > 0; g-- {
		if !f.isGroupHidden(g) {
			break
		}
		lastGroup--
	}

	for g, group := range f.groups {
		// determine the first non-skippable field.
		var firstField int
		for _, field := range group.fields {
			if !field.Skip() || len(group.fields) == 1 {
				break
			}
			firstField++
		}

		// determine the last non-skippable field.
		var lastField int
		for i := len(group.fields) - 1; i > 0; i-- {
			lastField = i
			if !group.fields[i].Skip() || len(group.fields) == 1 {
				break
			}
		}

		for i, field := range group.fields {
			field.WithPosition(FieldPosition{
				Group:      g,
				Field:      i,
				FirstField: firstField,
				LastField:  lastField,
				FirstGroup: firstGroup,
				LastGroup:  lastGroup,
			})
		}
	}
	return f
}

// Errors returns the current groups' errors.
func (f *Form) Errors() []error {
	return f.groups[f.paginator.Page].Errors()
}

// Help returns the current groups' help.
func (f *Form) Help() help.Model {
	return f.groups[f.paginator.Page].help
}

// KeyBinds returns the current fields' keybinds.
func (f *Form) KeyBinds() []key.Binding {
	group := f.groups[f.paginator.Page]
	return group.fields[group.paginator.Page].KeyBinds()
}

// Get returns a result from the form.
func (f *Form) Get(key string) any {
	return f.results[key]
}

// GetString returns a result as a string from the form.
func (f *Form) GetString(key string) string {
	v, ok := f.results[key].(string)
	if !ok {
		return ""
	}
	return v
}

// GetInt returns a result as a string from the form.
func (f *Form) GetInt(key string) int {
	v, ok := f.results[key].(int)
	if !ok {
		return 0
	}
	return v
}

// GetBool returns a result as a string from the form.
func (f *Form) GetBool(key string) bool {
	v, ok := f.results[key].(bool)
	if !ok {
		return false
	}
	return v
}

// NextGroup moves the form to the next group.
func (f *Form) NextGroup() tea.Cmd {
	_, cmd := f.Update(nextGroup())
	return cmd
}

// PrevGroup moves the form to the next group.
func (f *Form) PrevGroup() tea.Cmd {
	_, cmd := f.Update(prevGroup())
	return cmd
}

// NextField moves the form to the next field.
func (f *Form) NextField() tea.Cmd {
	_, cmd := f.Update(NextField())
	return cmd
}

// NextField moves the form to the next field.
func (f *Form) PrevField() tea.Cmd {
	_, cmd := f.Update(PrevField())
	return cmd
}

// Init initializes the form.
func (f *Form) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(f.groups))
	for i, group := range f.groups {
		cmds[i] = group.Init()
	}

	if f.isGroupHidden(f.paginator.Page) {
		cmds = append(cmds, nextGroup)
	}

	return tea.Batch(cmds...)
}

// Update updates the form.
func (f *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If the form is aborted or completed there's no need to update it.
	if f.State != StateNormal {
		return f, nil
	}

	page := f.paginator.Page
	group := f.groups[page]

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if f.width > 0 {
			break
		}
		for _, group := range f.groups {
			group.WithWidth(msg.Width)
		}
		if f.height > 0 {
			break
		}
		for _, group := range f.groups {
			if group.fullHeight() > msg.Height {
				group.WithHeight(msg.Height)
			}
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, f.keymap.Quit):
			f.aborted = true
			f.quitting = true
			f.State = StateAborted
			return f, f.cancelCmd
		}

	case nextFieldMsg:
		// Form is progressing to the next field, let's save the value of the current field.
		field := group.fields[group.paginator.Page]
		f.results[field.GetKey()] = field.GetValue()

	case nextGroupMsg:
		if len(group.Errors()) > 0 {
			return f, nil
		}

		submit := func() (tea.Model, tea.Cmd) {
			f.quitting = true
			f.State = StateCompleted
			return f, f.submitCmd
		}

		if f.paginator.OnLastPage() {
			return submit()
		}

		for i := f.paginator.Page + 1; i < f.paginator.TotalPages; i++ {
			if !f.isGroupHidden(i) {
				f.paginator.Page = i
				break
			}
			// all subsequent groups are hidden, so we must act as
			// if we were in the last one.
			if i == f.paginator.TotalPages-1 {
				return submit()
			}
		}
		return f, f.groups[f.paginator.Page].Init()

	case prevGroupMsg:
		if len(group.Errors()) > 0 {
			return f, nil
		}

		for i := f.paginator.Page - 1; i >= 0; i-- {
			if !f.isGroupHidden(i) {
				f.paginator.Page = i
				break
			}
		}

		return f, f.groups[f.paginator.Page].Init()
	}

	m, cmd := group.Update(msg)
	f.groups[page] = m.(*Group)

	// A user input a key, this could hide or show other groups,
	// let's update all of their positions.
	switch msg.(type) {
	case tea.KeyMsg:
		f.UpdateFieldPositions()
	}

	return f, cmd
}

func (f *Form) isGroupHidden(page int) bool {
	hide := f.groups[page].hide
	if hide == nil {
		return false
	}
	return hide()
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
	m, err := tea.NewProgram(f, f.teaOptions...).Run()
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
