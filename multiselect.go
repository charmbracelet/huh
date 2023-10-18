package huh

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"
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
	style            *MultiSelectStyle
	blurredStyle     MultiSelectStyle
	focusedStyle     MultiSelectStyle
}

// NewMultiSelect returns a new multi-select field.
func NewMultiSelect() *MultiSelect {
	f, b := DefaultMultiSelectStyles()
	return &MultiSelect{
		cursorPrefix:     "> ",
		selectedPrefix:   "[•] ",
		unselectedPrefix: "[ ] ",
		focusedStyle:     f,
		blurredStyle:     b,
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

// Focus focuses the multi-select field.
func (m *MultiSelect) Focus() tea.Cmd {
	m.style = &m.focusedStyle
	return nil
}

// Blur blurs the multi-select field.
func (m *MultiSelect) Blur() tea.Cmd {
	m.style = &m.blurredStyle
	return nil
}

// Init initializes the multi-select field.
func (m *MultiSelect) Init() tea.Cmd {
	m.style = &m.blurredStyle
	return nil
}

// Update updates the multi-select field.
func (m *MultiSelect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.cursor = max(m.cursor-1, 0)
		case "down", "j":
			m.cursor = min(m.cursor+1, len(m.options)-1)
		case " ", "x":
			m.options[m.cursor].selected = !m.options[m.cursor].selected
			if m.options[m.cursor].selected {
				*m.value = append(*m.value, m.options[m.cursor].name)
			} else {
				i := slices.Index(*m.value, m.options[m.cursor].name)
				*m.value = slices.Delete(*m.value, i, i+1)
			}
		case "shift+tab":
			return m, prevField
		case "tab", "enter":
			return m, nextField
		}
	}

	return m, nil
}

// View renders the multi-select field.
func (m *MultiSelect) View() string {
	var sb strings.Builder
	sb.WriteString(m.style.Title.Render(m.title) + "\n")
	c := m.style.Cursor.Render(m.cursorPrefix)
	for i, option := range m.options {
		if m.cursor == i {
			sb.WriteString(c)
		} else {
			sb.WriteString(strings.Repeat(" ", lipgloss.Width(c)))
		}

		if option.selected {
			sb.WriteString(m.style.SelectedPrefix.Render(m.selectedPrefix))
			sb.WriteString(m.style.Selected.Render(option.name))
		} else {
			sb.WriteString(m.style.UnselectedPrefix.Render(m.unselectedPrefix))
			sb.WriteString(m.style.Unselected.Render(option.name))
		}
		if i < len(m.options)-1 {
			sb.WriteString("\n")
		}
	}
	return m.style.Base.Render(sb.String())
}

// Run runs the multi-select field in accessible mode.
func (m *MultiSelect) Run() {
	fmt.Println(m.style.Title.Render(m.title))

	for i, option := range m.options {
		fmt.Printf("%d. %s\n", i+1, option.name)
	}

	fmt.Println("\nType 0 to finish.\n")

	var choice int
	for {
		choice = accessibility.PromptInt(0, len(m.options))
		if choice == 0 {
			break
		}
		m.options[choice-1].selected = !m.options[choice-1].selected
		if m.options[choice-1].selected {
			fmt.Println("Selected:", m.options[choice-1].name)
		} else {
			fmt.Println("Unselected:", m.options[choice-1].name)
		}
	}

	for _, option := range m.options {
		if option.selected {
			*m.value = append(*m.value, option.name)
		}
	}

	fmt.Printf("Selected: %s\n\n", strings.Join(*m.value, ", "))
}
