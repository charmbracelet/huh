package main

import (
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const (
	focusColor = "#2EF8BB"
	breakColor = "#FF5F87"
)

var (
	focusTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(focusColor)).MarginRight(1).SetString("Focus Mode")
	breakTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(breakColor)).MarginRight(1).SetString("Break Mode")
	pausedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(breakColor)).MarginRight(1).SetString("Continue?")
	helpStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginTop(2)
	sidebarStyle    = lipgloss.NewStyle().MarginLeft(3).Padding(1, 3).Border(lipgloss.RoundedBorder()).BorderForeground(helpStyle.GetForeground())
)

var baseTimerStyle = lipgloss.NewStyle().Padding(1, 2)

type mode int

const (
	Initial mode = iota
	Focusing
	Paused
	Breaking
)

type Model struct {
	form     *huh.Form
	quitting bool

	lastTick  time.Time
	startTime time.Time

	mode mode

	focusTime time.Duration
	breakTime time.Duration

	progress progress.Model
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

const tickInterval = time.Second / 2

type tickMsg time.Time

func tickCmd(t time.Time) tea.Msg {
	return tickMsg(t)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tickMsg:
		cmds = append(cmds, tea.Tick(tickInterval, tickCmd))
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			switch m.mode {
			case Focusing:
				m.mode = Paused
				m.startTime = time.Now()
				m.progress.FullColor = breakColor
			case Paused:
				m.mode = Breaking
				m.startTime = time.Now()
			case Breaking:
				m.quitting = true
				return m, tea.Quit
			}
		case "ctrl+c":
			m.quitting = true
			return m, tea.Interrupt
		default:
			if m.mode == Paused {
				m.mode = Breaking
				m.startTime = time.Now()
			}
		}
	}

	// Update form
	f, cmd := m.form.Update(msg)
	m.form = f.(*huh.Form)
	cmds = append(cmds, cmd)
	if m.form.State != huh.StateCompleted {
		return m, tea.Batch(cmds...)
	}

	// Update timer
	if m.startTime.IsZero() {
		m.startTime = time.Now()
		m.focusTime = m.form.Get("focus").(time.Duration)
		m.breakTime = m.form.Get("break").(time.Duration)
		m.mode = Focusing
		cmds = append(cmds, tea.Tick(tickInterval, tickCmd))
	}

	switch m.mode {
	case Focusing:
		if time.Now().After(m.startTime.Add(m.focusTime)) {
			m.mode = Paused
			m.startTime = time.Now()
			m.progress.FullColor = breakColor
		}
	case Breaking:
		if time.Now().After(m.startTime.Add(m.breakTime)) {
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	if m.form.State != huh.StateCompleted {
		return m.form.View()
	}

	var s strings.Builder

	elapsed := time.Now().Sub(m.startTime)
	var percent float64
	switch m.mode {
	case Focusing:
		percent = float64(elapsed) / float64(m.focusTime)
		s.WriteString(focusTitleStyle.String())
		s.WriteString(elapsed.Round(time.Second).String())
		s.WriteString("\n\n")
		s.WriteString(m.progress.ViewAs(percent))
		s.WriteString(helpStyle.Render("Press 'q' to skip"))
	case Paused:
		s.WriteString(pausedStyle.String())
		s.WriteString("\n\nFocus time is done, time to take a break.")
		s.WriteString(helpStyle.Render("press any key to continue.\n"))
	case Breaking:
		percent = float64(elapsed) / float64(m.breakTime)
		s.WriteString(breakTitleStyle.String())
		s.WriteString(elapsed.Round(time.Second).String())
		s.WriteString("\n\n")
		s.WriteString(m.progress.ViewAs(percent))
		s.WriteString(helpStyle.Render("press 'q' to quit"))
	}

	return baseTimerStyle.Render(s.String())
}

func NewModel() Model {
	theme := huh.ThemeCharm()
	theme.Focused.Base.Border(lipgloss.HiddenBorder())
	theme.Focused.Title.Foreground(lipgloss.Color(focusColor))
	theme.Focused.SelectSelector.Foreground(lipgloss.Color(focusColor))
	theme.Focused.SelectedOption.Foreground(lipgloss.Color("15"))
	theme.Focused.Option.Foreground(lipgloss.Color("7"))

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[time.Duration]().
				Title("Focus Time").
				Key("focus").
				Options(
					huh.NewOption("25 minutes", 25*time.Minute),
					huh.NewOption("30 minutes", 30*time.Minute),
					huh.NewOption("45 minutes", 45*time.Minute),
					huh.NewOption("1 hour", time.Hour),
				),
		),
		huh.NewGroup(
			huh.NewSelect[time.Duration]().
				Title("Break Time").
				Key("break").
				Options(
					huh.NewOption("5 minutes", 5*time.Minute),
					huh.NewOption("10 minutes", 10*time.Minute),
					huh.NewOption("15 minutes", 15*time.Minute),
					huh.NewOption("20 minutes", 20*time.Minute),
				),
		),
	).WithShowHelp(false).WithTheme(theme)

	progress := progress.New()
	progress.FullColor = focusColor
	progress.SetSpringOptions(1, 1)

	return Model{
		form:     form,
		progress: progress,
	}
}

func main() {
	m := NewModel()
	mm, err := tea.NewProgram(&m).Run()
	m = mm.(Model)
	if err != nil {
		log.Fatal(err)
	}
}
