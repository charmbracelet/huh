package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const maxWidth = 80

var (
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 2, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(indigo).
		Bold(true).
		Padding(0, 1, 0, 4)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		PaddingLeft(1).
		MarginTop(1).
		MarginRight(2)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

type state int

const (
	statusNormal state = iota
	stateDone
)

type Model struct {
	state  state
	lg     *lipgloss.Renderer
	styles *Styles
	form   *huh.Form
	width  int
}

func NewModel() Model {
	m := Model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("class").
				Options(huh.NewOptions("Warrior", "Mage", "Rogue")...).
				Title("Choose your class").
				Description("This will determine your department"),

			huh.NewSelect[string]().
				Key("level").
				Options(huh.NewOptions("1", "20", "9999")...).
				Title("Choose your level").
				Description("This will determine your benefits package"),

			huh.NewConfirm().
				Title("All done?").
				Affirmative("Yep").
				Negative("Wait, no"),
		),
	).WithHelp(false)
	return m
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = max(msg.Width, maxWidth) - m.styles.Base.GetHorizontalFrameSize()
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	return m, cmd
}

func (m Model) View() string {
	s := m.styles

	var class string
	if m.getFormValue("class") != "" {
		class = "Class: " + m.getFormValue("class")
	}

	// Form (left side)
	v := m.form.View()
	if m.form.State == huh.StateCompleted {
		v += s.Highlight.Render(fmt.Sprintf("You selected: Level %s, %s\n", m.form.Get("level"), m.form.Get("class")))
		v += s.Help.Render("\nctrl+c to quit\n")
	}
	form := m.lg.NewStyle().Margin(1, 2).Render(v)

	// Status (right side)
	var status string
	{
		var (
			buildInfo      = "(None)"
			role           string
			jobDescription string
			level          string
		)

		if m.getFormValue("level") != "" {
			level = "Level: " + m.getFormValue("level")
			role, jobDescription = m.getRole()
			role = "\n\n" + s.StatusHeader.Render("Projected Role") + "\n" + role
			jobDescription = "\n\n" + s.StatusHeader.Render("Duties") + "\n" + jobDescription
		}
		if m.getFormValue("class") != "" {
			buildInfo = fmt.Sprintf("%s\n%s", class, level)
		}

		const statusWidth = 28
		statusMarginLeft := m.width - statusWidth - lipgloss.Width(form) - s.Status.GetMarginRight()
		status = s.Status.Copy().
			Height(lipgloss.Height(form)).
			Width(statusWidth).
			MarginLeft(statusMarginLeft).
			Render(s.StatusHeader.Render("Current Build") + "\n" +
				buildInfo +
				role +
				jobDescription)
	}

	header := m.appBoundaryView("Charm Employment Application")
	body := lipgloss.JoinHorizontal(lipgloss.Top, form, status)
	footer := m.appBoundaryView("")

	return s.Base.Render(header + "\n" + body + "\n\n" + footer)
}

func (m Model) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(indigo),
	)
}

func (m Model) getFormValue(key string) string {
	s, ok := m.form.Get(key).(string)
	if !ok {
		return ""
	}
	return s
}

func (m Model) getRole() (string, string) {
	level := m.getFormValue("level")
	switch m.getFormValue("class") {
	case "Warrior":
		switch level {
		case "1":
			return "Tank Intern", "Assists with tank-related activities. Paid position."
		case "9999":
			return "Tank Manager", "Manages tanks and tank-related activities."
		default:
			return "Tank", "General tank. Does damage, takes damage. Responsible for tanking."
		}
	case "Mage":
		switch level {
		case "1":
			return "DPS Associate", "Finds DPS deals and passes them on to DPS Manager."
		case "9999":
			return "DPS Operating Officer", "Oversees all DPS activities."
		default:
			return "DPS", "Does damage and ideally does not take damage. Logs hours in JIRA."
		}
	case "Rogue":
		switch level {
		case "1":
			return "Stealth Junior Designer", "Designs rougue-like activities. Reports to Stealth Lead."
		case "9999":
			return "Stealth Lead", "Lead designer for all things stealth. Some travel required."
		default:
			return "Sneaky Person", "Sneaks around and does sneaky things. Reports to Stealth Lead."
		}
	default:
		return "", ""
	}
}

func main() {
	_, err := tea.NewProgram(NewModel()).Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}
}
