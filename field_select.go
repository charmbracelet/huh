package huh

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"
)

// Select is a form select field.
type Select[T comparable] struct {
	value    *T
	key      string
	viewport viewport.Model

	// customization
	title           string
	description     string
	options         []Option[T]
	filteredOptions []Option[T]
	height          int

	// error handling
	validate func(T) error
	err      error

	// state
	selected  int
	focused   bool
	filtering bool
	filter    textinput.Model

	// options
	inline     bool
	width      int
	accessible bool
	theme      *Theme
	keymap     SelectKeyMap
}

// NewSelect returns a new select field.
func NewSelect[T comparable]() *Select[T] {
	filter := textinput.New()
	filter.Prompt = "/"

	return &Select[T]{
		options:   []Option[T]{},
		value:     new(T),
		validate:  func(T) error { return nil },
		filtering: false,
		filter:    filter,
	}
}

// Value sets the value of the select field.
func (s *Select[T]) Value(value *T) *Select[T] {
	s.value = value
	s.selectValue(*value)
	return s
}

func (s *Select[T]) selectValue(value T) {
	for i, o := range s.options {
		if o.Value == value {
			s.selected = i
			break
		}
	}
}

// Key sets the key of the select field which can be used to retrieve the value
// after submission.
func (s *Select[T]) Key(key string) *Select[T] {
	s.key = key
	return s
}

// Title sets the title of the select field.
func (s *Select[T]) Title(title string) *Select[T] {
	s.title = title
	return s
}

// Description sets the description of the select field.
func (s *Select[T]) Description(description string) *Select[T] {
	s.description = description
	return s
}

// Options sets the options of the select field.
func (s *Select[T]) Options(options ...Option[T]) *Select[T] {
	if len(options) <= 0 {
		return s
	}
	s.options = options
	s.filteredOptions = options

	// Set the cursor to the existing value or the last selected option.
	for i, option := range options {
		if option.Value == *s.value {
			s.selected = i
			break
		} else if option.selected {
			s.selected = i
		}
	}

	s.updateViewportHeight()

	return s
}

// Inline sets whether the select input should be inline.
func (s *Select[T]) Inline(v bool) *Select[T] {
	s.inline = v
	if v {
		s.Height(1)
	}
	s.keymap.Left.SetEnabled(v)
	s.keymap.Right.SetEnabled(v)
	s.keymap.Up.SetEnabled(!v)
	s.keymap.Down.SetEnabled(!v)
	return s
}

// Height sets the height of the select field. If the number of options
// exceeds the height, the select field will become scrollable.
func (s *Select[T]) Height(height int) *Select[T] {
	s.height = height
	s.updateViewportHeight()
	return s
}

// Validate sets the validation function of the select field.
func (s *Select[T]) Validate(validate func(T) error) *Select[T] {
	s.validate = validate
	return s
}

// Error returns the error of the select field.
func (s *Select[T]) Error() error {
	return s.err
}

// Skip returns whether the select should be skipped or should be blocking.
func (*Select[T]) Skip() bool {
	return false
}

// Zoom returns whether the input should be zoomed.
func (*Select[T]) Zoom() bool {
	return false
}

// Focus focuses the select field.
func (s *Select[T]) Focus() tea.Cmd {
	s.focused = true
	return nil
}

// Blur blurs the select field.
func (s *Select[T]) Blur() tea.Cmd {
	value := *s.value
	if s.inline {
		s.clearFilter()
		s.selectValue(value)
	}
	s.focused = false
	s.err = s.validate(value)
	return nil
}

// KeyBinds returns the help keybindings for the select field.
func (s *Select[T]) KeyBinds() []key.Binding {
	return []key.Binding{
		s.keymap.Up,
		s.keymap.Down,
		s.keymap.Left,
		s.keymap.Right,
		s.keymap.Filter,
		s.keymap.SetFilter,
		s.keymap.ClearFilter,
		s.keymap.Prev,
		s.keymap.Next,
		s.keymap.Submit,
	}
}

// Init initializes the select field.
func (s *Select[T]) Init() tea.Cmd {
	return nil
}

// Update updates the select field.
func (s *Select[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	s.updateViewportHeight()

	var cmd tea.Cmd
	if s.filtering {
		s.filter, cmd = s.filter.Update(msg)

		// Keep the selected item in view.
		if s.selected < s.viewport.YOffset || s.selected >= s.viewport.YOffset+s.viewport.Height {
			s.viewport.SetYOffset(s.selected)
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		s.err = nil
		switch {
		case key.Matches(msg, s.keymap.Filter):
			s.setFiltering(true)
			return s, s.filter.Focus()
		case key.Matches(msg, s.keymap.SetFilter):
			if len(s.filteredOptions) <= 0 {
				s.filter.SetValue("")
				s.filteredOptions = s.options
			}
			s.setFiltering(false)
		case key.Matches(msg, s.keymap.ClearFilter):
			s.clearFilter()
		case key.Matches(msg, s.keymap.Up, s.keymap.Left):
			// When filtering we should ignore j/k keybindings
			//
			// XXX: Currently, the below check doesn't account for keymap
			// changes. When making this fix it's worth considering ignoring
			// whether to ignore all up/down keybindings as ignoring a-zA-Z0-9
			// may not be enough when international keyboards are considered.
			if s.filtering && (msg.String() == "k" || msg.String() == "h") {
				break
			}
			s.selected = max(s.selected-1, 0)
			if s.selected < s.viewport.YOffset {
				s.viewport.SetYOffset(s.selected)
			}
		case key.Matches(msg, s.keymap.GotoTop):
			if s.filtering {
				break
			}
			s.selected = 0
			s.viewport.GotoTop()
		case key.Matches(msg, s.keymap.GotoBottom):
			if s.filtering {
				break
			}
			s.selected = len(s.filteredOptions) - 1
			s.viewport.GotoBottom()
		case key.Matches(msg, s.keymap.HalfPageUp):
			s.selected = max(s.selected-s.viewport.Height/2, 0)
			s.viewport.HalfViewUp()
		case key.Matches(msg, s.keymap.HalfPageDown):
			s.selected = min(s.selected+s.viewport.Height/2, len(s.filteredOptions)-1)
			s.viewport.HalfViewDown()
		case key.Matches(msg, s.keymap.Down, s.keymap.Right):
			// When filtering we should ignore j/k keybindings
			//
			// XXX: See note in the previous case match.
			if s.filtering && (msg.String() == "j" || msg.String() == "l") {
				break
			}
			s.selected = min(s.selected+1, len(s.filteredOptions)-1)
			if s.selected >= s.viewport.YOffset+s.viewport.Height {
				s.viewport.LineDown(1)
			}
		case key.Matches(msg, s.keymap.Prev):
			if s.selected >= len(s.filteredOptions) {
				break
			}
			value := s.filteredOptions[s.selected].Value
			s.err = s.validate(value)
			if s.err != nil {
				return s, nil
			}
			*s.value = value
			return s, PrevField
		case key.Matches(msg, s.keymap.Next, s.keymap.Submit):
			if s.selected >= len(s.filteredOptions) {
				break
			}
			value := s.filteredOptions[s.selected].Value
			s.setFiltering(false)
			s.err = s.validate(value)
			if s.err != nil {
				return s, nil
			}
			*s.value = value
			return s, NextField
		}

		if s.filtering {
			s.filteredOptions = s.options
			if s.filter.Value() != "" {
				s.filteredOptions = nil
				for _, option := range s.options {
					if s.filterFunc(option.Key) {
						s.filteredOptions = append(s.filteredOptions, option)
					}
				}
			}
			if len(s.filteredOptions) > 0 {
				s.selected = min(s.selected, len(s.filteredOptions)-1)
				s.viewport.SetYOffset(clamp(s.selected, 0, len(s.filteredOptions)-s.viewport.Height))
			}
		}
	}

	return s, cmd
}

// updateViewportHeight updates the viewport size according to the Height setting
// on this select field.
func (s *Select[T]) updateViewportHeight() {
	// If no height is set size the viewport to the number of options.
	if s.height <= 0 {
		s.viewport.Height = len(s.options)
		return
	}

	const minHeight = 1
	s.viewport.Height = max(minHeight, s.height-
		lipgloss.Height(s.titleView())-
		lipgloss.Height(s.descriptionView()))
}

func (s *Select[T]) activeStyles() *FieldStyles {
	theme := s.theme
	if theme == nil {
		theme = ThemeCharm()
	}
	if s.focused {
		return &theme.Focused
	}
	return &theme.Blurred
}

func (s *Select[T]) titleView() string {
	if s.title == "" {
		return ""
	}
	var (
		styles = s.activeStyles()
		sb     = strings.Builder{}
	)
	if s.filtering {
		sb.WriteString(styles.Title.Render(s.filter.View()))
	} else if s.filter.Value() != "" && !s.inline {
		sb.WriteString(styles.Title.Render(s.title) + styles.Description.Render("/"+s.filter.Value()))
	} else {
		sb.WriteString(styles.Title.Render(s.title))
	}
	if s.err != nil {
		sb.WriteString(styles.ErrorIndicator.String())
	}
	return sb.String()
}

func (s *Select[T]) descriptionView() string {
	return s.activeStyles().Description.Render(s.description)
}

func (s *Select[T]) choicesView() string {
	var (
		styles = s.activeStyles()
		c      = styles.SelectSelector.String()
		sb     strings.Builder
	)

	if s.inline {
		sb.WriteString(styles.PrevIndicator.Faint(s.selected <= 0).String())
		if len(s.filteredOptions) > 0 {
			sb.WriteString(styles.SelectedOption.Render(s.filteredOptions[s.selected].Key))
		} else {
			sb.WriteString(styles.TextInput.Placeholder.Render("No matches"))
		}
		sb.WriteString(styles.NextIndicator.Faint(s.selected == len(s.filteredOptions)-1).String())
		return sb.String()
	}

	for i, option := range s.filteredOptions {
		if s.selected == i {
			sb.WriteString(c + styles.SelectedOption.Render(option.Key))
		} else {
			sb.WriteString(strings.Repeat(" ", lipgloss.Width(c)) + styles.Option.Render(option.Key))
		}
		if i < len(s.options)-1 {
			sb.WriteString("\n")
		}
	}

	for i := len(s.filteredOptions); i < len(s.options)-1; i++ {
		sb.WriteString("\n")
	}

	return sb.String()
}

// View renders the select field.
func (s *Select[T]) View() string {
	styles := s.activeStyles()
	s.viewport.SetContent(s.choicesView())

	var sb strings.Builder
	if s.title != "" {
		sb.WriteString(s.titleView())
		if !s.inline {
			sb.WriteString("\n")
		}
	}
	if s.description != "" {
		sb.WriteString(s.descriptionView())
		if !s.inline {
			sb.WriteString("\n")
		}
	}
	sb.WriteString(s.viewport.View())
	return styles.Base.Render(sb.String())
}

// clearFilter clears the value of the filter.
func (s *Select[T]) clearFilter() {
	s.filter.SetValue("")
	s.filteredOptions = s.options
	s.setFiltering(false)
}

// setFiltering sets the filter of the select field.
func (s *Select[T]) setFiltering(filtering bool) {
	if s.inline && filtering {
		s.filter.Width = lipgloss.Width(s.titleView()) - 1 - 1
	}
	s.filtering = filtering
	s.keymap.SetFilter.SetEnabled(filtering)
	s.keymap.Filter.SetEnabled(!filtering)
	s.keymap.ClearFilter.SetEnabled(!filtering && s.filter.Value() != "")
}

// filterFunc returns true if the option matches the filter.
func (s *Select[T]) filterFunc(option string) bool {
	// XXX: remove diacritics or allow customization of filter function.
	return strings.Contains(strings.ToLower(option), strings.ToLower(s.filter.Value()))
}

// Run runs the select field.
func (s *Select[T]) Run() error {
	if s.accessible {
		return s.runAccessible()
	}
	return Run(s)
}

// runAccessible runs an accessible select field.
func (s *Select[T]) runAccessible() error {
	var sb strings.Builder
	styles := s.activeStyles()

	sb.WriteString(styles.Title.Render(s.title) + "\n")

	for i, option := range s.options {
		sb.WriteString(fmt.Sprintf("%d. %s", i+1, option.Key))
		sb.WriteString("\n")
	}

	fmt.Println(sb.String())

	for {
		choice := accessibility.PromptInt("Choose: ", 1, len(s.options))
		option := s.options[choice-1]
		if err := s.validate(option.Value); err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(styles.SelectedOption.Render("Chose: " + option.Key + "\n"))
		*s.value = option.Value
		break
	}

	return nil
}

// WithTheme sets the theme of the select field.
func (s *Select[T]) WithTheme(theme *Theme) Field {
	if s.theme != nil {
		return s
	}
	s.theme = theme
	s.filter.Cursor.Style = s.theme.Focused.TextInput.Cursor
	s.filter.PromptStyle = s.theme.Focused.TextInput.Prompt
	s.updateViewportHeight()
	return s
}

// WithKeyMap sets the keymap on a select field.
func (s *Select[T]) WithKeyMap(k *KeyMap) Field {
	s.keymap = k.Select
	s.keymap.Left.SetEnabled(s.inline)
	s.keymap.Right.SetEnabled(s.inline)
	s.keymap.Up.SetEnabled(!s.inline)
	s.keymap.Down.SetEnabled(!s.inline)
	return s
}

// WithAccessible sets the accessible mode of the select field.
func (s *Select[T]) WithAccessible(accessible bool) Field {
	s.accessible = accessible
	return s
}

// WithWidth sets the width of the select field.
func (s *Select[T]) WithWidth(width int) Field {
	s.width = width
	return s
}

// WithHeight sets the height of the select field.
func (s *Select[T]) WithHeight(height int) Field {
	return s.Height(height)
}

// WithPosition sets the position of the select field.
func (s *Select[T]) WithPosition(p FieldPosition) Field {
	if s.filtering {
		return s
	}
	s.keymap.Prev.SetEnabled(!p.IsFirst())
	s.keymap.Next.SetEnabled(!p.IsLast())
	s.keymap.Submit.SetEnabled(p.IsLast())
	return s
}

// GetKey returns the key of the field.
func (s *Select[T]) GetKey() string {
	return s.key
}

// GetValue returns the value of the field.
func (s *Select[T]) GetValue() any {
	return *s.value
}
