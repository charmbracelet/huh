// Package modal provides a reusable overlay-dialog component for Bubble Tea v2
// applications, built on top of the lipgloss v2 compositor.
//
// A Modal wraps any [Content] (typically a [huh.Form]) and is rendered as a
// layer stacked on top of the parent application. While a modal is open the
// parent should route messages to it exclusively, which gives the modal
// exclusive focus until it resolves.
//
// Lifecycle:
//
//  1. Construct a Modal with [New] (or [NewForm] for a huh.Form).
//  2. The parent owns the modal as an optional field (nil = closed).
//  3. On every parent Update, if the field is non-nil forward the message
//     to Modal.Update only, then watch the returned command for a [ResolvedMsg].
//  4. On [ResolvedMsg], inspect Confirmed and Value, then clear the field.
//  5. On every parent View, if the field is non-nil call [Modal.Render] to
//     composite it over the background string.
package modal

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Content is anything that can live inside a [Modal]. It mirrors the standard
// Bubble Tea model lifecycle plus a [Content.Done] check that lets the modal
// know when to resolve.
type Content interface {
	Init() tea.Cmd
	Update(tea.Msg) (Content, tea.Cmd)
	View() string
	// Done reports whether the content has resolved. When done is true the
	// modal will emit a [ResolvedMsg] carrying confirmed and value. A
	// cancelled resolution should return (true, false, nil).
	Done() (done bool, confirmed bool, value any)
}

// ResolvedMsg is dispatched once when the modal's [Content] reports done.
// Parents should clear their modal field on receipt and react to Confirmed
// and Value as appropriate.
type ResolvedMsg struct {
	ID        string
	Confirmed bool
	Value     any
}

// Modal wraps a [Content] for rendering as a centered overlay layer.
type Modal struct {
	id       string
	content  Content
	style    lipgloss.Style
	z        int
	resolved bool
}

// Option configures a [Modal] at construction time.
type Option func(*Modal)

// WithStyle wraps the content with the given lipgloss style (typically a
// bordered, padded box). Defaults to a rounded border with single-cell
// padding.
func WithStyle(s lipgloss.Style) Option {
	return func(m *Modal) { m.style = s }
}

// WithZ sets the z-index of the modal layer. Defaults to 10. Higher values
// render on top of lower ones when multiple modals are stacked.
func WithZ(z int) Option {
	return func(m *Modal) { m.z = z }
}

// defaultStyle is a rounded-border dialog box. Callers can override via
// [WithStyle].
var defaultStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#874BFD")).
	Padding(1, 2)

// New constructs a Modal for the given Content. The id is used to identify
// the modal in [ResolvedMsg] (and as the layer ID for hit testing).
func New(id string, content Content, opts ...Option) *Modal {
	m := &Modal{
		id:      id,
		content: content,
		style:   defaultStyle,
		z:       10,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// ID returns the modal's identifier.
func (m *Modal) ID() string { return m.id }

// Init returns the underlying content's initial command.
func (m *Modal) Init() tea.Cmd {
	return m.content.Init()
}

// Update forwards the message to the content. If the content reports it is
// done, Update also emits a [ResolvedMsg]. Once resolved the modal becomes
// inert — further Update calls are no-ops.
func (m *Modal) Update(msg tea.Msg) (*Modal, tea.Cmd) {
	if m.resolved {
		return m, nil
	}
	var cmd tea.Cmd
	m.content, cmd = m.content.Update(msg)
	done, confirmed, value := m.content.Done()
	if !done {
		return m, cmd
	}
	m.resolved = true
	id := m.id
	resolveCmd := func() tea.Msg {
		return ResolvedMsg{ID: id, Confirmed: confirmed, Value: value}
	}
	if cmd == nil {
		return m, resolveCmd
	}
	return m, tea.Batch(cmd, resolveCmd)
}

// View renders the styled content as a plain string (no compositing). Use
// [Modal.Layer] or [Modal.Render] for placement over a background.
func (m *Modal) View() string {
	return m.style.Render(m.content.View())
}

// Layer returns the modal as a positioned lipgloss layer, centered within a
// canvas of size (parentW, parentH). Use this when assembling your own layer
// tree; for the simple "background + one modal" case prefer [Modal.Render].
func (m *Modal) Layer(parentW, parentH int) *lipgloss.Layer {
	view := m.View()
	w := lipgloss.Width(view)
	h := lipgloss.Height(view)
	x := (parentW - w) / 2
	y := (parentH - h) / 2
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	return lipgloss.NewLayer(view).ID(m.id).X(x).Y(y).Z(m.z)
}

// Render composites the modal over background and returns the final string
// ready to drop into tea.View.Content. parentW and parentH must match the
// dimensions of background so centering is correct.
func (m *Modal) Render(background string, parentW, parentH int) string {
	root := lipgloss.NewLayer(background).ID("modal-background")
	root.AddLayers(m.Layer(parentW, parentH))
	return lipgloss.NewCompositor(root).Render()
}
