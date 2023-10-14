package huh

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/ordered"
	"golang.org/x/exp/slices"
)

type option struct {
	name     string
	selected bool
}

// MultiSelect is a form multi-select field.
type MultiSelect struct {
	title            string
	required         bool
	filterable       bool
	limit            int
	cursor           int
	cursorPrefix     string
	selectedPrefix   string
	unselectedPrefix string
	selected         []int
	options          []option
	value            *[]string
}

// NewMultiSelect returns a new multi-select field.
func NewMultiSelect() *MultiSelect {
	return &MultiSelect{
		cursorPrefix:     " > ",
		selectedPrefix:   "[•]",
		unselectedPrefix: "[ ]",
	}
}

// Value sets the value of the multi-select field.
func (m *MultiSelect) Value(value *[]string) *MultiSelect {
	m.value = value
	return m
}

// Title sets the title of the multi-select field.
func (m *MultiSelect) Title(title string) *MultiSelect {
	m.title = title
	return m
}

// Required sets the multi-select field as required.
func (m *MultiSelect) Required(required bool) *MultiSelect {
	m.required = required
	return m
}

// Options sets the options of the multi-select field.
func (m *MultiSelect) Options(options ...string) *MultiSelect {
	for _, o := range options {
		m.options = append(m.options, option{o, false})
	}
	return m
}

// Filterable sets the multi-select field as filterable.
func (m *MultiSelect) Filterable(filterable bool) *MultiSelect {
	m.filterable = filterable
	return m
}

// Cursor sets the cursor of the multi-select field.
func (m *MultiSelect) Cursor(cursor string) *MultiSelect {
	m.cursorPrefix = cursor
	return m
}

// Limit sets the limit of the multi-select field.
func (m *MultiSelect) Limit(limit int) *MultiSelect {
	m.limit = limit
	return m
}

// Init initializes the multi-select field.
func (m *MultiSelect) Init() tea.Cmd {
	return nil
}

// Update updates the multi-select field.
func (m *MultiSelect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.cursor = ordered.Max(m.cursor-1, 0)
		case "down", "j":
			m.cursor = ordered.Min(m.cursor+1, len(m.options)-1)
		case " ", "x":
			m.options[m.cursor].selected = !m.options[m.cursor].selected
			if m.options[m.cursor].selected {
				*m.value = append(*m.value, m.options[m.cursor].name)
			} else {
				i := slices.Index(*m.value, m.options[m.cursor].name)
				*m.value = slices.Delete(*m.value, i, i+1)
			}
		case "tab", "enter":
			return m, nextField
		}
	}

	return m, nil
}

// View renders the multi-select field.
func (m *MultiSelect) View() string {
	var sb strings.Builder
	sb.WriteString(m.title + "\n")
	for i, option := range m.options {
		if m.cursor == i {
			sb.WriteString(m.cursorPrefix)
		} else {
			sb.WriteString(strings.Repeat(" ", lipgloss.Width(m.cursorPrefix)))
		}

		if option.selected {
			sb.WriteString(m.selectedPrefix)
		} else {
			sb.WriteString(m.unselectedPrefix)
		}

		sb.WriteString(option.name)

		if i < len(m.options)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
