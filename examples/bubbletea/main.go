package main

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

const maxWidth = 80

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style

	Red, Indigo, Green color.Color
}

func NewStyles(hasDarkBg bool) *Styles {
	var (
		s         = Styles{}
		lightDark = lipgloss.LightDark(hasDarkBg)
	)

	s.Red = lightDark(lipgloss.Color("#FE5F86"), lipgloss.Color("#FE5F86"))
	s.Indigo = lightDark(lipgloss.Color("#5A56E0"), lipgloss.Color("#7571F9"))
	s.Green = lightDark(lipgloss.Color("#02BA84"), lipgloss.Color("#02BF87"))
	s.Base = lipgloss.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lipgloss.NewStyle().
		Foreground(s.Indigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lipgloss.NewStyle().
		Foreground(s.Green).
		Bold(true)
	s.Highlight = lipgloss.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(s.Red)
	s.Help = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

type state int

const (
	statusNormal state = iota
	stateDone
)

type Model struct {
	state     state
	styles    func(bool) *Styles
	form      *huh.Form
	hasDarkBg bool
	width     int
}

func NewModel() Model {
	m := Model{
		width:  maxWidth,
		styles: NewStyles,
	}

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
				Key("done").
				Title("All done?").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("Welp, finish up then")
					}
					return nil
				}).
				Affirmative("Yep").
				Negative("Wait, no"),
		),
	).
		WithWidth(45).
		WithShowHelp(false).
		WithShowErrors(false)
	return m
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	styles := m.styles(m.hasDarkBg)
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		m.hasDarkBg = msg.IsDark()
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - styles.Base.GetHorizontalFrameSize()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Interrupt
		case "esc", "q":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		// Quit when the form is done.
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	s := m.styles(m.hasDarkBg)

	switch m.form.State {
	case huh.StateCompleted:
		title, role := m.getRole()
		title = s.Highlight.Render(title)
		var b strings.Builder
		fmt.Fprintf(&b, "Congratulations, you’re Charm’s newest\n%s!\n\n", title)
		fmt.Fprintf(&b, "Your job description is as follows:\n\n%s\n\nPlease proceed to HR immediately.", role)
		return tea.NewView(s.Status.Margin(0, 1).Padding(1, 2).Width(48).Render(b.String()) + "\n\n")
	default:

		var class string
		if m.form.GetString("class") != "" {
			class = "Class: " + m.form.GetString("class")
		}

		// Form (left side)
		v := strings.TrimSuffix(m.form.View(), "\n\n")
		form := lipgloss.NewStyle().Margin(1, 0).Render(v)

		// Status (right side)
		var status string
		{
			var (
				buildInfo      = "(None)"
				role           string
				jobDescription string
				level          string
			)

			if m.form.GetString("level") != "" {
				level = "Level: " + m.form.GetString("level")
				role, jobDescription = m.getRole()
				role = "\n\n" + s.StatusHeader.Render("Projected Role") + "\n" + role
				jobDescription = "\n\n" + s.StatusHeader.Render("Duties") + "\n" + jobDescription
			}
			if m.form.GetString("class") != "" {
				buildInfo = fmt.Sprintf("%s\n%s", class, level)
			}

			const statusWidth = 28
			statusMarginLeft := m.width - statusWidth - lipgloss.Width(form) - s.Status.GetMarginRight()
			status = s.Status.
				Height(lipgloss.Height(form)).
				Width(statusWidth).
				MarginLeft(statusMarginLeft).
				Render(s.StatusHeader.Render("Current Build") + "\n" +
					buildInfo +
					role +
					jobDescription)
		}

		errors := m.form.Errors()
		header := m.appBoundaryView("Charm Employment Application")
		if len(errors) > 0 {
			header = m.appErrorBoundaryView(m.errorView())
		}
		body := lipgloss.JoinHorizontal(lipgloss.Left, form, status)

		footer := m.appBoundaryView(m.form.Help().ShortHelpView(m.form.KeyBinds()))
		if len(errors) > 0 {
			footer = m.appErrorBoundaryView("")
		}

		return tea.NewView(s.Base.Render(header + "\n" + body + "\n\n" + footer))
	}
}

func (m Model) errorView() string {
	var s string
	for _, err := range m.form.Errors() {
		s += err.Error()
	}
	return s
}

func (m Model) appBoundaryView(text string) string {
	s := m.styles(m.hasDarkBg)
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		s.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceStyle(lipgloss.NewStyle().Foreground(s.Indigo)),
	)
}

func (m Model) appErrorBoundaryView(text string) string {
	s := m.styles(m.hasDarkBg)
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		s.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceStyle(lipgloss.NewStyle().Foreground(s.Red)),
	)
}

func (m Model) getRole() (string, string) {
	level := m.form.GetString("level")
	switch m.form.GetString("class") {
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
			return "Stealth Junior Designer", "Designs rogue-like activities. Reports to Stealth Lead."
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
