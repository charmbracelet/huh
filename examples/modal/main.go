package main

import (
	"fmt"
	"os"
	"strconv"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/huh/v2/modal"
	"charm.land/lipgloss/v2"
)

const wizardID = "character-wizard"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#874BFD")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	tableBorderStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240"))

	backgroundStyle = lipgloss.NewStyle().Padding(2, 4)
)

type character struct {
	class     string
	level     string
	eyeColor  string
	hairColor string
}

type model struct {
	width, height int
	modal         *modal.Modal
	characters    []character
	table         table.Model
	lastCancelled bool
}

func newModel() model {
	cols := []table.Column{
		{Title: "#", Width: 3},
		{Title: "Class", Width: 9},
		{Title: "Level", Width: 6},
		{Title: "Eyes", Width: 7},
		{Title: "Hair", Width: 7},
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithHeight(8),
		table.WithWidth(40),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return model{table: t}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case modal.ResolvedMsg:
		if msg.ID == wizardID {
			confirmed := msg.Confirmed
			var saved *character
			if confirmed {
				form := msg.Value.(*huh.Form)
				if form.GetBool("ok") {
					c := character{
						class:     form.GetString("class"),
						level:     form.GetString("level"),
						eyeColor:  form.GetString("eyeColor"),
						hairColor: form.GetString("hairColor"),
					}
					saved = &c
				} else {
					confirmed = false
				}
			}
			if saved != nil {
				m.characters = append(m.characters, *saved)
				m.refreshRows()
				m.lastCancelled = false
			} else {
				m.lastCancelled = !confirmed
			}
			m.modal = nil
		}
		return m, nil
	}

	// While a modal is open, route messages to it exclusively.
	if m.modal != nil {
		var cmd tea.Cmd
		m.modal, cmd = m.modal.Update(msg)
		return m, cmd
	}

	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "n":
			m.modal = newWizard()
			return m, m.modal.Init()
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *model) refreshRows() {
	rows := make([]table.Row, len(m.characters))
	for i, c := range m.characters {
		rows[i] = table.Row{strconv.Itoa(i + 1), c.class, c.level, c.eyeColor, c.hairColor}
	}
	m.table.SetRows(rows)
}

func (m model) View() tea.View {
	bg := m.backgroundContent()
	bg = lipgloss.Place(m.width, m.height, lipgloss.Top, lipgloss.Left, bg)

	v := tea.NewView(bg)
	v.AltScreen = true
	if m.modal != nil {
		v.SetContent(m.modal.Render(bg, m.width, m.height))
	}
	return v
}

func (m model) backgroundContent() string {
	body := titleStyle.Render(" Character Roster ") + "\n\n"
	body += tableBorderStyle.Render(m.table.View()) + "\n"

	if m.lastCancelled {
		body += "\n" + helpStyle.Render("(last wizard cancelled — nothing saved)")
	}

	body += "\n\n" + helpStyle.Render("n — new character • ↑/↓ — navigate • q — quit")
	return backgroundStyle.Render(body)
}

func newWizard() *modal.Modal {
	colors := huh.NewOptions("Brown", "Black", "Green")

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("class").
				Title("Choose your class").
				Options(huh.NewOptions("Warrior", "Mage", "Rogue")...),

			huh.NewSelect[string]().
				Key("level").
				Title("Choose your level").
				Options(huh.NewOptions("1", "20", "9999")...),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("eyeColor").
				Title("Choose your eye color").
				Options(colors...),

			huh.NewSelect[string]().
				Key("hairColor").
				Title("Choose your hair color").
				Options(colors...),

			huh.NewConfirm().
				Key("ok").
				Title("Save this character?").
				Affirmative("Save").
				Negative("Cancel"),
		),
	).
		WithWidth(45).
		WithShowHelp(true)

	return modal.NewForm(wizardID, form)
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
