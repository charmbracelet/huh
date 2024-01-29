package huh

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"
)

// File is a form file file field.
type File struct {
	value  *string
	key    string
	picker filepicker.Model

	// state
	focused bool

	// customization
	title       string
	description string

	// error handling
	validate func(string) error
	err      error

	// options
	width      int
	accessible bool
	theme      *Theme
	keymap     FileKeyMap
}

const defaultHeight = 5

// NewFile returns a new file field.
func NewFile() *File {
	fp := filepicker.New()
	fp.ShowPermissions = false
	fp.ShowSize = false
	fp.Height = defaultHeight
	fp.AutoHeight = false

	cmd := fp.Init()
	if cmd != nil {
		fp, _ = fp.Update(cmd())
	}

	return &File{
		value:    new(string),
		validate: func(string) error { return nil },
		picker:   fp,
		theme:    ThemeCharm(),
	}
}

// CurrentDirectory sets the directory of the file field.
func (f *File) CurrentDirectory(directory string) *File {
	f.picker.CurrentDirectory = directory
	return f
}

// ShowHidden sets whether to show hidden files.
func (f *File) ShowHidden(v bool) *File {
	f.picker.ShowHidden = v
	return f
}

// Value sets the value of the file field.
func (f *File) Value(value *string) *File {
	f.value = value
	return f
}

// Key sets the key of the file field which can be used to retrieve the value
// after submission.
func (f *File) Key(key string) *File {
	f.key = key
	return f
}

// Title sets the title of the file field.
func (f *File) Title(title string) *File {
	f.title = title
	return f
}

// Description sets the description of the file field.
func (f *File) Description(description string) *File {
	f.description = description
	return f
}

// Height sets the height of the file field. If the number of options
// exceeds the height, the file field will become scrollable.
func (f *File) AllowedTypes(types []string) *File {
	f.picker.AllowedTypes = types
	return f
}

// Height sets the height of the file field. If the number of options
// exceeds the height, the file field will become scrollable.
func (f *File) Height(height int) *File {
	f.picker.Height = height
	f.picker.AutoHeight = false
	return f
}

// Validate sets the validation function of the file field.
func (f *File) Validate(validate func(string) error) *File {
	f.validate = validate
	return f
}

// Error returns the error of the file field.
func (f *File) Error() error {
	return f.err
}

// Skip returns whether the file should be skipped or should be blocking.
func (*File) Skip() bool {
	return false
}

// Focus focuses the file field.
func (f *File) Focus() tea.Cmd {
	f.focused = true
	return f.picker.Init()
}

// Blur blurs the file field.
func (f *File) Blur() tea.Cmd {
	f.focused = false
	f.err = f.validate(*f.value)
	return nil
}

// KeyBinds returns the help keybindings for the file field.
func (f *File) KeyBinds() []key.Binding {
	return []key.Binding{f.keymap.Up, f.keymap.Down, f.keymap.Prev, f.keymap.Next, f.keymap.Submit}
}

// Init initializes the file field.
func (f *File) Init() tea.Cmd {
	return f.picker.Init()
}

// Update updates the file field.
func (f *File) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	f.err = nil

	var cmd tea.Cmd
	f.picker, cmd = f.picker.Update(msg)
	didSelect, file := f.picker.DidSelectFile(msg)
	if didSelect {
		*f.value = file
		return f, nextField
	}
	didSelect, file = f.picker.DidSelectDisabledFile(msg)
	if didSelect {
		f.err = errors.New("cannot select " + filepath.Base(file))
		return f, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, f.keymap.Next):
			return f, nextField
		case key.Matches(msg, f.keymap.Prev):
			return f, prevField
		}
	}

	return f, cmd
}

// View renders the file field.
func (f *File) View() string {
	styles := f.theme.Blurred
	if f.focused {
		styles = f.theme.Focused
	}
	var sb strings.Builder
	if f.title != "" {
		sb.WriteString(styles.Title.Render(f.title) + "\n")
	}
	if f.description != "" {
		sb.WriteString(styles.Description.Render(f.description) + "\n")
	}
	sb.WriteString(strings.TrimSuffix(f.picker.View(), "\n"))
	return styles.Base.Render(sb.String())
}

// Run runs the file field.
func (f *File) Run() error {
	if f.accessible {
		return f.runAccessible()
	}
	return Run(f)
}

// runAccessible runs an accessible file field.
func (f *File) runAccessible() error {
	fmt.Println(f.theme.Blurred.Base.Render(f.theme.Focused.Title.Render(f.title)))
	fmt.Println()

	validateFile := func(s string) error {
		// is the string a file?
		if _, err := os.Open(s); err != nil {
			return errors.New("not a file")
		}

		// is it one of the allowed types?
		valid := false
		for _, ext := range f.picker.AllowedTypes {
			if strings.HasSuffix(s, ext) {
				valid = true
				break
			}
		}
		if !valid {
			return errors.New("cannot select: " + s)
		}

		// does it pass user validation?
		return f.validate(s)
	}

	*f.value = accessibility.PromptString("File: ", validateFile)
	fmt.Println(f.theme.Focused.SelectedOption.Render("File: " + *f.value + "\n"))
	return nil
}

// WithTheme sets the theme of the file field.
func (f *File) WithTheme(theme *Theme) Field {
	f.theme = theme

	// TODO: add specific themes
	f.picker.Styles = filepicker.Styles{
		DisabledCursor:   lipgloss.Style{},
		Cursor:           theme.Focused.TextInput.Prompt,
		Symlink:          lipgloss.NewStyle(),
		Directory:        theme.Focused.Title,
		File:             lipgloss.NewStyle(),
		DisabledFile:     theme.Focused.Description,
		Permission:       theme.Focused.Description,
		Selected:         theme.Focused.SelectedOption,
		DisabledSelected: theme.Focused.Description,
		FileSize:         theme.Focused.Description.Copy().Width(7).Align(lipgloss.Right).Inline(true),
		EmptyDirectory:   theme.Focused.Description.Copy().SetString("No files found."),
	}

	return f
}

// WithKeyMap sets the keymap on a file field.
func (f *File) WithKeyMap(k *KeyMap) Field {
	f.keymap = k.File
	f.picker.KeyMap = filepicker.KeyMap{
		GoToTop:  k.File.GoToTop,
		GoToLast: k.File.GoToLast,
		Down:     k.File.Down,
		Up:       k.File.Up,
		PageUp:   k.File.PageUp,
		PageDown: k.File.PageDown,
		Back:     k.File.Back,
		Open:     k.File.Open,
		Select:   k.File.Select,
	}
	return f
}

// WithAccessible sets the accessible mode of the file field.
func (f *File) WithAccessible(accessible bool) Field {
	f.accessible = accessible
	return f
}

// WithWidth sets the width of the file field.
func (f *File) WithWidth(width int) Field {
	f.width = width
	return f
}

// WithHeight sets the height of the file field.
func (f *File) WithHeight(height int) Field {
	return f.Height(height)
}

// WithPosition sets the position of the file field.
func (f *File) WithPosition(p FieldPosition) Field {
	f.keymap.Prev.SetEnabled(!p.IsFirst())
	f.keymap.Next.SetEnabled(!p.IsLast())
	f.keymap.Submit.SetEnabled(p.IsLast())
	return f
}

// GetKey returns the key of the field.
func (f *File) GetKey() string {
	return f.key
}

// GetValue returns the value of the field.
func (f *File) GetValue() any {
	return *f.value
}
