package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var highlight = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render
var help = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render

type Model struct {
	form *huh.Form
}

func NewModel() Model {
	var m Model
	f := huh.NewForm(
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
		),
	)

	m.form = f
	return m
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
	v := "Charm Employment Application\n\n" + m.form.View()
	if m.form.State == huh.StateCompleted {
		v += highlight(fmt.Sprintf("You selected: Level %s, %s\n", m.form.Get("level"), m.form.Get("class")))
		v += help("\nctrl+c to quit\n")
	}
	return lipgloss.NewStyle().Margin(1, 2).Render(v)
}

func main() {
	_, err := tea.NewProgram(NewModel()).Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}
}
