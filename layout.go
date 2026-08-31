package huh

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// A Layout is responsible for laying out groups in a form.
type Layout interface {
	View(f *Form) string
	GroupWidth(f *Form, g *Group, w int) int
}

// LayoutDefault is the default layout shows a single group at a time.
var LayoutDefault Layout = &layoutDefault{}

// LayoutStack is a layout stacks all groups on top of each other.
var LayoutStack Layout = &layoutStack{}

// LayoutColumns layout distributes groups in even columns.
func LayoutColumns(columns int) Layout {
	return &layoutColumns{columns: columns}
}

// LayoutGrid layout distributes groups in a grid.
func LayoutGrid(rows int, columns int) Layout {
	return &layoutGrid{rows: rows, columns: columns}
}

// visibleFormGroups returns all groups that are not currently hidden, in
// order, along with the position of the currently selected group *within
// that filtered slice*. If the selected group is itself hidden, pos is -1.
func visibleFormGroups(f *Form) (groups []*Group, pos int) {
	pos = -1
	selectedIdx := f.selector.Index()

	f.selector.Range(func(i int, group *Group) bool {
		if f.isGroupHidden(group) {
			return true
		}
		if i == selectedIdx {
			pos = len(groups)
		}
		groups = append(groups, group)
		return true
	})

	return groups, pos
}

type layoutDefault struct{}

func (l *layoutDefault) View(f *Form) string {
	return f.selector.Selected().View()
}

func (l *layoutDefault) GroupWidth(_ *Form, _ *Group, w int) int {
	return w
}

type layoutColumns struct {
	columns int
}

func (l *layoutColumns) visibleGroups(f *Form) []*Group {
	groups, pos := visibleFormGroups(f)
	if len(groups) == 0 {
		return nil
	}
	if pos < 0 {
		pos = 0
	}

	segmentIndex := pos / l.columns
	start := segmentIndex * l.columns
	end := min(start+l.columns, len(groups))

	return groups[start:end]
}

func (l *layoutColumns) View(f *Form) string {
	groups := l.visibleGroups(f)
	if len(groups) == 0 {
		return ""
	}

	columns := make([]string, 0, len(groups))
	for _, group := range groups {
		columns = append(columns, group.Content())
	}

	header := f.selector.Selected().Header()
	footer := f.selector.Selected().Footer()

	return strings.Join([]string{
		header,
		lipgloss.JoinHorizontal(lipgloss.Left, columns...),
		footer,
	}, "\n")
}

func (l *layoutColumns) GroupWidth(_ *Form, _ *Group, w int) int {
	return w / l.columns
}

type layoutStack struct{}

func (l *layoutStack) View(f *Form) string {
	groups, _ := visibleFormGroups(f)

	var columns []string
	for _, group := range groups {
		columns = append(columns, group.Content(), "")
	}

	if footer := f.selector.Selected().Footer(); footer != "" {
		columns = append(columns, footer)
	}
	return strings.Join(columns, "\n")
}

func (l *layoutStack) GroupWidth(_ *Form, _ *Group, w int) int {
	return w
}

type layoutGrid struct {
	rows, columns int
}

func (l *layoutGrid) visibleGroups(f *Form) [][]*Group {
	groups, pos := visibleFormGroups(f)
	if len(groups) == 0 {
		return nil
	}
	if pos < 0 {
		pos = 0
	}

	total := l.rows * l.columns
	segmentIndex := pos / total
	start := segmentIndex * total
	end := min(start+total, len(groups))
	page := groups[start:end]

	grid := make([][]*Group, l.rows)
	for i := 0; i < l.rows; i++ {
		startRow := i * l.columns
		endRow := startRow + l.columns
		if startRow >= len(page) {
			break
		}
		if endRow > len(page) {
			endRow = len(page)
		}
		grid[i] = page[startRow:endRow]
	}
	return grid
}

func (l *layoutGrid) View(f *Form) string {
	grid := l.visibleGroups(f)
	if len(grid) == 0 {
		return ""
	}

	rows := make([]string, 0, len(grid))
	for _, row := range grid {
		var columns []string
		for _, group := range row {
			columns = append(columns, group.Content())
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Left, columns...), "")
	}
	footer := f.selector.Selected().Footer()

	return strings.Join(append(rows, footer), "\n")
}

func (l *layoutGrid) GroupWidth(_ *Form, _ *Group, w int) int {
	return w / l.columns
}
