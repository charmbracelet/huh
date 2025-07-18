package huh

import (
	"cmp"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/internal/accessibility"
	"github.com/charmbracelet/lipgloss"
)

const (
	minHeight     = 1
	defaultHeight = 10
)

// Select is a select field.
//
// A select field is a field that allows the user to select from a list of
// options. The options can be provided statically or dynamically using Options
// or OptionsFunc. The options can be filtered using "/" and navigation is done
// using j/k, up/down, or ctrl+n/ctrl+p keys.
type Select[T comparable] struct {
	id       int
	accessor Accessor[T]
	key      string

	viewport viewport.Model

	title           Eval[string]
	description     Eval[string]
	options         Eval[[]Option[T]]
	filteredOptions []Option[T]

	validate func(T) error
	err      error

	selected  int
	focused   bool
	filtering bool
	filter    textinput.Model
	spinner   spinner.Model

	inline     bool
	width      int
	height     int
	accessible bool // Deprecated: use RunAccessible instead.
	theme      *Theme
	keymap     SelectKeyMap
}

// NewSelect creates a new select field.
//
// A select field is a field that allows the user to select from a list of
// options. The options can be provided statically or dynamically using Options
// or OptionsFunc. The options can be filtered using "/" and navigation is done
// using j/k, up/down, or ctrl+n/ctrl+p keys.
func NewSelect[T comparable]() *Select[T] {
	filter := textinput.New()
	filter.Prompt = "/"

	s := spinner.New(spinner.WithSpinner(spinner.Line))

	return &Select[T]{
		accessor:    &EmbeddedAccessor[T]{},
		validate:    func(T) error { return nil },
		filtering:   false,
		filter:      filter,
		options:     Eval[[]Option[T]]{cache: make(map[uint64][]Option[T])},
		title:       Eval[string]{cache: make(map[uint64]string)},
		description: Eval[string]{cache: make(map[uint64]string)},
		spinner:     s,
	}
}

// Value sets the value of the select field.
func (s *Select[T]) Value(value *T) *Select[T] {
	return s.Accessor(NewPointerAccessor(value))
}

// Accessor sets the accessor of the select field.
func (s *Select[T]) Accessor(accessor Accessor[T]) *Select[T] {
	s.accessor = accessor
	s.selectValue(s.accessor.Get())
	s.updateValue()
	return s
}

func (s *Select[T]) selectValue(value T) {
	for i, o := range s.options.val {
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
//
// This title will be static, for dynamic titles use `TitleFunc`.
func (s *Select[T]) Title(title string) *Select[T] {
	s.title.val = title
	s.title.fn = nil
	return s
}

// TitleFunc sets the title func of the select field.
//
// This TitleFunc will be re-evaluated when the binding of the TitleFunc
// changes. This when you want to display dynamic content and update the title
// when another part of your form changes.
//
// See README#Dynamic for more usage information.
func (s *Select[T]) TitleFunc(f func() string, bindings any) *Select[T] {
	s.title.fn = f
	s.title.bindings = bindings
	return s
}

// Filtering sets the filtering state of the select field.
func (s *Select[T]) Filtering(filtering bool) *Select[T] {
	s.filtering = filtering
	s.filter.Focus()
	return s
}

// Description sets the description of the select field.
//
// This description will be static, for dynamic descriptions use `DescriptionFunc`.
func (s *Select[T]) Description(description string) *Select[T] {
	s.description.val = description
	return s
}

// DescriptionFunc sets the description func of the select field.
//
// This DescriptionFunc will be re-evaluated when the binding of the
// DescriptionFunc changes. This is useful when you want to display dynamic
// content and update the description when another part of your form changes.
//
// See README#Dynamic for more usage information.
func (s *Select[T]) DescriptionFunc(f func() string, bindings any) *Select[T] {
	s.description.fn = f
	s.description.bindings = bindings
	return s
}

// Options sets the options of the select field.
//
// This is what your user will select from.
//
// Title
// Description
//
//	-> Option 1
//	   Option 2
//	   Option 3
//
// These options will be static, for dynamic options use `OptionsFunc`.
func (s *Select[T]) Options(options ...Option[T]) *Select[T] {
	if len(options) <= 0 {
		return s
	}
	s.options.val = options
	s.filteredOptions = options

	s.selectOption()

	s.updateViewportHeight()
	s.updateValue()

	return s
}

func (s *Select[T]) selectOption() {
	// Set the cursor to the existing value or the last selected option.
	for i, option := range s.options.val {
		if option.Value == s.accessor.Get() {
			s.selected = i
			break
		}
		if option.selected {
			s.selected = i
			break
		}
	}
	s.viewport.YOffset = s.selected
}

// OptionsFunc sets the options func of the select field.
//
// This OptionsFunc will be re-evaluated when the binding of the OptionsFunc
// changes. This is useful when you want to display dynamic content and update
// the options when another part of your form changes.
//
// For example, changing the state / provinces, based on the selected country.
//
//	   huh.NewSelect[string]().
//		    Options(huh.NewOptions("United States", "Canada", "Mexico")...).
//		    Value(&country).
//		    Title("Country").
//		    Height(5),
//
//		huh.NewSelect[string]().
//		  Title("State / Province"). // This can also be made dynamic with `TitleFunc`.
//		  OptionsFunc(func() []huh.Option[string] {
//		    s := states[country]
//		    time.Sleep(1000 * time.Millisecond)
//		    return huh.NewOptions(s...)
//		}, &country),
//
// See examples/dynamic/dynamic-country/main.go for the full example.
func (s *Select[T]) OptionsFunc(f func() []Option[T], bindings any) *Select[T] {
	s.options.fn = f
	s.options.bindings = bindings
	// If there is no height set, we should attach a static height since these
	// options are possibly dynamic.
	if s.height <= 0 {
		s.height = defaultHeight
		s.updateViewportHeight()
	}
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

// Height sets the height of the select field. If the number of options exceeds
// the height, the select field will become scrollable.
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
func (s *Select[T]) Error() error { return s.err }

// Skip returns whether the select should be skipped or should be blocking.
func (*Select[T]) Skip() bool { return false }

// Zoom returns whether the input should be zoomed.
func (*Select[T]) Zoom() bool { return false }

// Focus focuses the select field.
func (s *Select[T]) Focus() tea.Cmd {
	s.focused = true
	return nil
}

// Blur blurs the select field.
func (s *Select[T]) Blur() tea.Cmd {
	value := s.accessor.Get()
	if s.inline {
		s.clearFilter()
		s.selectValue(value)
	}
	s.focused = false
	s.err = s.validate(value)
	return nil
}

// Hovered returns the value of the option under the cursor, and a bool
// indicating whether one was found. If there are no visible options, returns
// a zero-valued T and false.
func (s *Select[T]) Hovered() (T, bool) {
	if len(s.filteredOptions) == 0 || s.selected >= len(s.filteredOptions) {
		var zero T
		return zero, false
	}
	return s.filteredOptions[s.selected].Value, true
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
	}

	switch msg := msg.(type) {
	case updateFieldMsg:
		var cmds []tea.Cmd
		if ok, hash := s.title.shouldUpdate(); ok {
			s.title.bindingsHash = hash
			if !s.title.loadFromCache() {
				s.title.loading = true
				cmds = append(cmds, func() tea.Msg {
					return updateTitleMsg{id: s.id, title: s.title.fn(), hash: hash}
				})
			}
		}
		if ok, hash := s.description.shouldUpdate(); ok {
			s.description.bindingsHash = hash
			if !s.description.loadFromCache() {
				s.description.loading = true
				cmds = append(cmds, func() tea.Msg {
					return updateDescriptionMsg{id: s.id, description: s.description.fn(), hash: hash}
				})
			}
		}
		if ok, hash := s.options.shouldUpdate(); ok {
			s.clearFilter()
			s.options.bindingsHash = hash
			if s.options.loadFromCache() {
				s.filteredOptions = s.options.val
				s.selected = clamp(s.selected, 0, len(s.options.val)-1)
			} else {
				s.options.loading = true
				s.options.loadingStart = time.Now()
				cmds = append(cmds, func() tea.Msg {
					return updateOptionsMsg[T]{id: s.id, hash: hash, options: s.options.fn()}
				}, s.spinner.Tick)
			}
		}
		return s, tea.Batch(cmds...)

	case spinner.TickMsg:
		if !s.options.loading {
			break
		}
		s.spinner, cmd = s.spinner.Update(msg)
		return s, cmd

	case updateTitleMsg:
		if msg.id == s.id && msg.hash == s.title.bindingsHash {
			s.title.update(msg.title)
		}
	case updateDescriptionMsg:
		if msg.id == s.id && msg.hash == s.description.bindingsHash {
			s.description.update(msg.description)
		}
	case updateOptionsMsg[T]:
		if msg.id == s.id && msg.hash == s.options.bindingsHash {
			s.options.update(msg.options)
			s.selectOption()

			// since we're updating the options, we need to update the selected
			// cursor position and filteredOptions.
			s.selected = clamp(s.selected, 0, len(msg.options)-1)
			s.filteredOptions = msg.options
			s.updateValue()
		}
	case tea.KeyMsg:
		s.err = nil
		switch {
		case key.Matches(msg, s.keymap.Filter):
			s.setFiltering(true)
			return s, s.filter.Focus()
		case key.Matches(msg, s.keymap.SetFilter):
			if len(s.filteredOptions) <= 0 {
				s.filter.SetValue("")
				s.filteredOptions = s.options.val
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
			s.selected = s.selected - 1
			if s.selected < 0 {
				s.selected = len(s.filteredOptions) - 1
				s.viewport.GotoBottom()
			}
			if s.selected < s.viewport.YOffset {
				s.viewport.SetYOffset(s.selected)
			}
			s.updateValue()
		case key.Matches(msg, s.keymap.GotoTop):
			if s.filtering {
				break
			}
			s.selected = 0
			s.viewport.GotoTop()
			s.updateValue()
		case key.Matches(msg, s.keymap.GotoBottom):
			if s.filtering {
				break
			}
			s.selected = len(s.filteredOptions) - 1
			s.viewport.GotoBottom()
		case key.Matches(msg, s.keymap.HalfPageUp):
			s.selected = max(s.selected-s.viewport.Height/2, 0)
			s.viewport.HalfPageUp()
			s.updateValue()
		case key.Matches(msg, s.keymap.HalfPageDown):
			s.selected = min(s.selected+s.viewport.Height/2, len(s.filteredOptions)-1)
			s.viewport.HalfPageDown()
			s.updateValue()
		case key.Matches(msg, s.keymap.Down, s.keymap.Right):
			// When filtering we should ignore j/k keybindings
			//
			// XXX: See note in the previous case match.
			if s.filtering && (msg.String() == "j" || msg.String() == "l") {
				break
			}
			s.selected = s.selected + 1
			if s.selected > len(s.filteredOptions)-1 {
				s.selected = 0
				s.viewport.GotoTop()
			}
			if s.selected >= s.viewport.YOffset+s.viewport.Height {
				s.viewport.ScrollDown(1)
			}
			s.updateValue()
		case key.Matches(msg, s.keymap.Prev):
			if s.selected >= len(s.filteredOptions) {
				break
			}
			s.updateValue()
			s.err = s.validate(s.accessor.Get())
			if s.err != nil {
				return s, nil
			}
			s.updateValue()
			return s, PrevField
		case key.Matches(msg, s.keymap.Next, s.keymap.Submit):
			if s.selected >= len(s.filteredOptions) {
				break
			}
			s.setFiltering(false)
			s.updateValue()
			s.err = s.validate(s.accessor.Get())
			if s.err != nil {
				return s, nil
			}
			s.updateValue()
			return s, NextField
		}

		if s.filtering {
			s.filteredOptions = s.options.val
			if s.filter.Value() != "" {
				s.filteredOptions = nil
				for _, option := range s.options.val {
					if s.filterFunc(option.Key) {
						s.filteredOptions = append(s.filteredOptions, option)
					}
				}
			}
			if len(s.filteredOptions) > 0 {
				s.selected = min(s.selected, len(s.filteredOptions)-1)
			}
		}

		_, offset, height := s.optionsView()
		if offset > -1 && height > 0 && (offset < s.viewport.YOffset || height+offset >= s.viewport.YOffset+s.viewport.Height) {
			s.viewport.SetYOffset(offset)
		}
	}

	return s, cmd
}

func (s *Select[T]) updateValue() {
	if s.selected < len(s.filteredOptions) && s.selected >= 0 {
		s.accessor.Set(s.filteredOptions[s.selected].Value)
	}
}

// updateViewportHeight updates the viewport size according to the Height setting
// on this select field.
func (s *Select[T]) updateViewportHeight() {
	// If no height is set size the viewport to the number of options.
	if s.height <= 0 {
		v, _, _ := s.optionsView()
		s.viewport.Height = lipgloss.Height(v)
		return
	}

	offset := 0
	if ss := s.titleView(); ss != "" {
		offset += lipgloss.Height(ss)
	}
	if ss := s.descriptionView(); ss != "" {
		offset += lipgloss.Height(ss)
	}

	s.viewport.Height = max(minHeight, s.height-offset)
	s.viewport.YOffset = s.selected
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
	var (
		styles   = s.activeStyles()
		sb       = strings.Builder{}
		maxWidth = s.width - styles.Base.GetHorizontalFrameSize()
	)
	if s.filtering {
		sb.WriteString(s.filter.View())
	} else if s.filter.Value() != "" && !s.inline {
		sb.WriteString(styles.Description.Render("/" + s.filter.Value()))
	} else {
		sb.WriteString(styles.Title.Render(wrap(s.title.val, maxWidth)))
	}
	if s.err != nil {
		sb.WriteString(styles.ErrorIndicator.String())
	}
	return sb.String()
}

func (s *Select[T]) descriptionView() string {
	if s.description.val == "" {
		return ""
	}
	maxWidth := s.width - s.activeStyles().Base.GetHorizontalFrameSize()
	return s.activeStyles().Description.Render(wrap(s.description.val, maxWidth))
}

func (s *Select[T]) optionsView() (string, int, int) {
	var (
		styles = s.activeStyles()
		sb     strings.Builder
	)

	if s.options.loading && time.Since(s.options.loadingStart) > spinnerShowThreshold {
		s.spinner.Style = s.activeStyles().MultiSelectSelector.UnsetString()
		sb.WriteString(s.spinner.View() + " Loading...")
		return sb.String(), -1, 1
	}

	if s.inline {
		option := styles.TextInput.Placeholder.Render("No matches")
		if len(s.filteredOptions) > 0 {
			option = styles.SelectedOption.Render(s.filteredOptions[s.selected].Key)
		}
		return lipgloss.NewStyle().
				Width(s.width).
				Render(lipgloss.JoinHorizontal(
					lipgloss.Left,
					styles.PrevIndicator.Faint(s.selected <= 0).String(),
					option,
					styles.NextIndicator.Faint(s.selected == len(s.filteredOptions)-1).String(),
				)),
			-1, 1
	}

	var cursorOffset int
	var cursorHeight int
	for i, option := range s.filteredOptions {
		selected := s.selected == i
		line := s.renderOption(option, selected)
		if i < s.selected {
			cursorOffset += lipgloss.Height(line)
		}
		if selected {
			cursorHeight = lipgloss.Height(line)
		}

		sb.WriteString(line)
		if i < len(s.options.val)-1 {
			sb.WriteString("\n")
		}
	}

	for i := len(s.filteredOptions); i < len(s.options.val)-1; i++ {
		sb.WriteString("\n")
	}

	return sb.String(), cursorOffset, cursorHeight
}

func (s *Select[T]) renderOption(option Option[T], selected bool) string {
	var (
		styles   = s.activeStyles()
		cursor   = styles.SelectSelector.String()
		cursorW  = lipgloss.Width(cursor)
		maxWidth = s.width - s.activeStyles().Base.GetHorizontalFrameSize() - cursorW
	)

	key := wrap(option.Key, maxWidth)

	if selected {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			cursor,
			styles.SelectedOption.Render(key),
		)
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		strings.Repeat(" ", cursorW),
		styles.UnselectedOption.Render(key),
	)
}

// View renders the select field.
func (s *Select[T]) View() string {
	styles := s.activeStyles()
	vpc, _, _ := s.optionsView()
	s.viewport.SetContent(vpc)

	var parts []string
	if s.title.val != "" || s.title.fn != nil {
		parts = append(parts, s.titleView())
	}
	if s.description.val != "" || s.description.fn != nil {
		parts = append(parts, s.descriptionView())
	}
	parts = append(parts, s.viewport.View())
	return styles.Base.Width(s.width).Height(s.height).
		Render(strings.Join(parts, "\n"))
}

// clearFilter clears the value of the filter.
func (s *Select[T]) clearFilter() {
	s.filter.SetValue("")
	s.filteredOptions = s.options.val
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
	if s.accessible { // TODO: remove in a future release.
		return s.RunAccessible(os.Stdout, os.Stdin)
	}
	return Run(s)
}

// RunAccessible runs an accessible select field.
func (s *Select[T]) RunAccessible(w io.Writer, r io.Reader) error {
	styles := s.activeStyles()
	_, _ = fmt.Fprintln(w, styles.Title.
		PaddingRight(1).
		Render(cmp.Or(s.title.val, "Select:")))

	for i, option := range s.options.val {
		_, _ = fmt.Fprintf(w, "%d. %s\n", i+1, option.Key)
	}

	var defaultValue *int
	switch s.accessor.(type) {
	case *PointerAccessor[T]: // if its of this type, it means it has a default value
		s.selectOption() // make sure s.selected is set
		idx := s.selected + 1
		defaultValue = &idx
	}
	prompt := fmt.Sprintf("Enter a number between %d and %d: ", 1, len(s.options.val))
	if len(s.options.val) == 1 {
		prompt = "There is only one option available; enter the number 1:"
	}
	for {
		choice := accessibility.PromptInt(w, r, prompt, 1, len(s.options.val), defaultValue)
		option := s.options.val[choice-1]
		if err := s.validate(option.Value); err != nil {
			_, _ = fmt.Fprintln(w, err.Error())
			_, _ = fmt.Fprintln(w)
			continue
		}
		s.accessor.Set(option.Value)
		return nil
	}
}

// WithTheme sets the theme of the select field.
func (s *Select[T]) WithTheme(theme *Theme) Field {
	if s.theme != nil {
		return s
	}
	s.theme = theme
	s.filter.Cursor.Style = theme.Focused.TextInput.Cursor
	s.filter.Cursor.TextStyle = theme.Focused.TextInput.CursorText
	s.filter.PromptStyle = theme.Focused.TextInput.Prompt
	s.filter.TextStyle = theme.Focused.TextInput.Text
	s.filter.PlaceholderStyle = theme.Focused.TextInput.Placeholder
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
//
// Deprecated: you may now call [Select.RunAccessible] directly to run the
// field in accessible mode.
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
func (s *Select[T]) GetKey() string { return s.key }

// GetValue returns the value of the field.
func (s *Select[T]) GetValue() any {
	return s.accessor.Get()
}

// GetFiltering returns the filtering state of the field.
func (s *Select[T]) GetFiltering() bool {
	return s.filtering
}
