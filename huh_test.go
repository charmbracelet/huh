package huh

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/xpty"
)

var pretty = lipgloss.NewStyle().
	Width(60).
	Border(lipgloss.NormalBorder()).
	MarginTop(1).
	Padding(1, 3, 1, 2)

func TestForm(t *testing.T) {
	type Taco struct {
		Shell    string
		Base     string
		Toppings []string
	}

	type Order struct {
		Taco         Taco
		Name         string
		Instructions string
		Discount     bool
	}

	var taco Taco
	order := Order{Taco: taco}

	f := NewForm(
		NewGroup(
			NewSelect[string]().
				Options(NewOptions("Soft", "Hard")...).
				Title("Shell?").
				Description("Our tortillas are made fresh in-house every day.").
				Validate(func(t string) error {
					if t == "Hard" {
						return fmt.Errorf("we're out of hard shells, sorry")
					}
					return nil
				}).
				Value(&order.Taco.Shell),

			NewSelect[string]().
				Options(NewOptions("Chicken", "Beef", "Fish", "Beans")...).
				Value(&order.Taco.Base).
				Title("Base"),
		),

		// Prompt for toppings and special instructions.
		// The customer can ask for up to 4 toppings.
		NewGroup(
			NewMultiSelect[string]().
				Title("Toppings").
				Description("Choose up to 4.").
				Options(
					NewOption("Lettuce", "lettuce").Selected(true),
					NewOption("Tomatoes", "tomatoes").Selected(true),
					NewOption("Corn", "corn"),
					NewOption("Salsa", "salsa"),
					NewOption("Sour Cream", "sour cream"),
					NewOption("Cheese", "cheese"),
				).
				Validate(func(t []string) error {
					if len(t) <= 0 {
						return fmt.Errorf("at least one topping is required")
					}
					return nil
				}).
				Value(&order.Taco.Toppings).
				Filterable(true).
				Limit(4),
		),

		// Gather final details for the order.
		NewGroup(
			NewInput().
				Value(&order.Name).
				Title("What's your name?").
				Placeholder("Margaret Thatcher").
				Description("For when your order is ready."),

			NewText().
				Value(&order.Instructions).
				Placeholder("Just put it in the mailbox please").
				Title("Special Instructions").
				Description("Anything we should know?").
				CharLimit(400),

			NewConfirm().
				Title("Would you like 15% off?").
				Value(&order.Discount).
				Affirmative("Yes!").
				Negative("No."),
		),
	)

	f.Update(f.Init())

	view := ansi.Strip(f.View())

	//
	//  ┃ Shell?
	//  ┃ Our tortillas are made fresh in-house every day.
	//  ┃ > Soft
	//  ┃   Hard
	//
	//    Base
	//    > Chicken
	//      Beef
	//      Fish
	//      Beans
	//
	//   ↑ up • ↓ down • / filter • enter select
	//

	if !strings.Contains(view, "┃ Shell?") {
		t.Log(pretty.Render(view))
		t.Error("Expected form to contain Shell? title")
	}

	if !strings.Contains(view, "Our tortillas are made fresh in-house every day.") {
		t.Log(pretty.Render(view))
		t.Error("Expected form to contain tortilla description")
	}

	if !strings.Contains(view, "Base") {
		t.Log(pretty.Render(view))
		t.Error("Expected form to contain Base title")
	}

	// Attempt to select hard shell and retrieve error.
	m, _ := f.Update(keys('j'))
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	view = ansi.Strip(m.View())

	if !strings.Contains(view, "* we're out of hard shells, sorry") {
		t.Log(pretty.Render(view))
		t.Error("Expected form to show out of hard shells error")
	}

	m, _ = m.Update(keys('k'))

	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = batchUpdate(m, cmd)

	view = ansi.Strip(m.View())

	if !strings.Contains(view, "┃ > Chicken") {
		t.Log(pretty.Render(view))
		t.Fatal("Expected form to continue to base group")
	}

	// batchMsg + nextGroup
	m, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = batchUpdate(m, cmd)
	view = ansi.Strip(m.View())

	//
	// ┃ Toppings
	// ┃ Choose up to 4.
	// ┃ > ✓ Lettuce
	// ┃   ✓ Tomatoes
	// ┃   • Corn
	// ┃   • Salsa
	// ┃   • Sour Cream
	// ┃   • Cheese
	//
	//  x toggle • ↑ up • ↓ down • enter confirm • shift+tab back
	//
	if !strings.Contains(view, "Toppings") {
		t.Log(pretty.Render(view))
		t.Fatal("Expected form to show toppings group")
	}

	if !strings.Contains(view, "Choose up to 4.") {
		t.Log(pretty.Render(view))
		t.Error("Expected form to show toppings description")
	}

	if !strings.Contains(view, "> ✓ Lettuce ") {
		t.Log(pretty.Render(view))
		t.Error("Expected form to preselect lettuce")
	}

	if !strings.Contains(view, "  ✓ Tomatoes") {
		t.Log(pretty.Render(view))
		t.Error("Expected form to preselect tomatoes")
	}

	m, _ = m.Update(keys('j'))
	m, _ = m.Update(keys('j'))
	view = ansi.Strip(m.View())

	if !strings.Contains(view, "> • Corn") {
		t.Log(pretty.Render(view))
		t.Error("Expected form to change selection to corn")
	}

	m, _ = m.Update(keys('x'))
	view = ansi.Strip(m.View())

	if !strings.Contains(view, "> ✓ Corn") {
		t.Log(pretty.Render(view))
		t.Error("Expected form to change selection to corn")
	}

	m = batchUpdate(m.Update(tea.KeyMsg{Type: tea.KeyEnter}))
	view = ansi.Strip(m.View())

	if !strings.Contains(view, "What's your name?") {
		t.Log(pretty.Render(view))
		t.Error("Expected form to prompt for name")
	}

	if !strings.Contains(view, "Special Instructions") {
		t.Log(pretty.Render(view))
		t.Error("Expected form to prompt for special instructions")
	}

	if !strings.Contains(view, "Would you like 15% off?") {
		t.Log(pretty.Render(view))
		t.Error("Expected form to prompt for discount")
	}

	//
	// ┃ What's your name?
	// ┃ For when your order is ready.
	// ┃ > Margaret Thatcher
	//
	//    Special Instructions
	//    Anything we should know?
	//    Just put it in the mailbox please
	//
	//    Would you like 15% off?
	//
	//      Yes!     No.
	//
	//   enter next • shift+tab back
	//
	m.Update(keys('G', 'l', 'e', 'n'))
	view = ansi.Strip(m.View())
	if !strings.Contains(view, "Glen") {
		t.Log(pretty.Render(view))
		t.Error("Expected form to accept user input")
	}

	if order.Taco.Shell != "Soft" {
		t.Error("Expected order shell to be Soft")
	}

	if order.Taco.Base != "Chicken" {
		t.Error("Expected order shell to be Chicken")
	}

	if len(order.Taco.Toppings) != 3 {
		t.Error("Expected order to have 3 toppings")
	}

	if order.Name != "Glen" {
		t.Error("Expected order name to be Glen")
	}

	// TODO: Finish and submit form.
}

func TestInput(t *testing.T) {
	field := NewInput()
	f := NewForm(NewGroup(field))
	f.Update(f.Init())

	view := ansi.Strip(f.View())

	if !strings.Contains(view, ">") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain prompt.")
	}

	// Type Huh in the form.
	m, _ := f.Update(keys('H', 'u', 'h'))
	f = m.(*Form)
	view = ansi.Strip(f.View())

	if !strings.Contains(view, "Huh") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain Huh.")
	}

	if !strings.Contains(view, "enter submit") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain help.")
	}

	if field.GetValue() != "Huh" {
		t.Error("Expected field value to be Huh")
	}
}

func TestInlineInput(t *testing.T) {
	field := NewInput().
		Title("Input ").
		Prompt(": ").
		Description("Description").
		Inline(true)

	f := NewForm(NewGroup(field)).WithWidth(40)
	f.Update(f.Init())

	view := ansi.Strip(f.View())

	if !strings.Contains(view, "┃ Input Description:") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain inline input.")
	}

	// Type Huh in the form.
	m, _ := f.Update(keys('H', 'u', 'h'))
	f = m.(*Form)
	view = ansi.Strip(f.View())

	if !strings.Contains(view, "Huh") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain Huh.")
	}

	if !strings.Contains(view, "enter submit") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain help.")
	}

	if !strings.Contains(view, "┃ Input Description: Huh") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain help.")
	}

	if field.GetValue() != "Huh" {
		t.Error("Expected field value to be Huh")
	}
}

func TestText(t *testing.T) {
	field := NewText()
	f := NewForm(NewGroup(field))
	f.Update(f.Init())

	// Type Huh in the form.
	m, _ := f.Update(keys('H', 'u', 'h'))
	f = m.(*Form)
	view := ansi.Strip(f.View())

	if !strings.Contains(view, "Huh") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain Huh.")
	}

	if !strings.Contains(view, "alt+enter / ctrl+j new line • ctrl+e open editor • enter submit") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain help.")
	}

	if field.GetValue() != "Huh" {
		t.Error("Expected field value to be Huh")
	}
}

func TestTextExternalEditorHidden(t *testing.T) {
	field := NewText().ExternalEditor(false)
	f := NewForm(NewGroup(field))
	f.Update(f.Init())

	// Type Huh in the form.
	m, _ := f.Update(keys('H', 'u', 'h'))
	f = m.(*Form)
	view := ansi.Strip(f.View())

	if !strings.Contains(view, "Huh") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain Huh.")
	}

	if strings.Contains(view, "ctrl+e open editor") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain help without ctrl+e.")
	}

	if field.GetValue() != "Huh" {
		t.Error("Expected field value to be Huh")
	}
}

func TestConfirm(t *testing.T) {
	field := NewConfirm().Title("Are you sure?")
	f := NewForm(NewGroup(field))
	f.Update(f.Init())

	// Type Huh in the form.
	m, _ := f.Update(keys('H'))
	f = m.(*Form)
	view := ansi.Strip(f.View())

	if !strings.Contains(view, "Yes") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain Yes.")
	}

	if !strings.Contains(view, "No") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain No.")
	}

	if !strings.Contains(view, "Are you sure?") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain Are you sure?.")
	}

	if !strings.Contains(view, "←/→ toggle • enter submit") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain help.")
	}

	if field.GetValue() != false {
		t.Error("Expected field value to be false")
	}

	// Toggle left
	f.Update(tea.KeyMsg{Type: tea.KeyLeft})

	if field.GetValue() != true {
		t.Error("Expected field value to be true")
	}

	// Toggle right
	f.Update(tea.KeyMsg{Type: tea.KeyRight})

	if field.GetValue() != false {
		t.Error("Expected field value to be false")
	}
}

func TestSelect(t *testing.T) {
	field := NewSelect[string]().
		Options(NewOptions(
			"Foo\nLine 2",
			"Bar\nLine 2",
			"Baz\nLine 2",
			"Ban\nLine 2",
		)...).
		Title("Which one?")
	f := NewForm(NewGroup(field)).WithHeight(5)
	f.Update(f.Init())

	view := ansi.Strip(f.View())

	if !strings.Contains(view, "Foo") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain Foo.")
	}

	if !strings.Contains(view, "Which one?") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain Which one?.")
	}

	if !strings.Contains(view, "> Foo") {
		t.Log(pretty.Render(view))
		t.Error("Expected cursor to be on Foo.")
	}

	// Move selection cursor down
	m, _ := f.Update(tea.KeyMsg{Type: tea.KeyDown})
	f = m.(*Form)

	view = ansi.Strip(f.View())

	if got, ok := field.Hovered(); !ok || got != "Bar\nLine 2" {
		t.Log(pretty.Render(view))
		t.Error("Expected cursor to be on Bar.")
	}

	if strings.Contains(view, "> Foo") {
		t.Log(pretty.Render(view))
		t.Error("Expected cursor to be on Bar.")
	}

	if !strings.Contains(view, "> Bar") {
		t.Log(pretty.Render(view))
		t.Error("Expected cursor to be on Bar.")
	}

	if !strings.Contains(view, "↑ up • ↓ down • / filter • enter submit") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain help.")
	}

	// Submit
	f.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if field.GetValue() != "Bar\nLine 2" {
		t.Error("Expected field value to be Bar")
	}
}

// doAllUpdates updates the form with the given command, then continues updating it with any resultant commands from the update until no more are returned.
func doAllUpdates(f *Form, cmd tea.Cmd) {
	if cmd == nil {
		return
	}
	var cmds []tea.Cmd
	switch msg := cmd().(type) {
	case tea.BatchMsg:
		for _, subcommand := range msg {
			doAllUpdates(f, subcommand)
		}
		return
	default:
		_, result := f.Update(msg)
		cmds = append(cmds, result)
	}
	doAllUpdates(f, tea.Batch(cmds...))
}

func TestSelectDynamic(t *testing.T) {
	trigger := "initial"

	field1 := NewSelect[string]().
		TitleFunc(func() string {
			return "field1 title " + trigger
		}, &trigger).
		DescriptionFunc(func() string {
			return "field1 desc " + trigger
		}, &trigger).
		OptionsFunc(func() []Option[string] {
			return []Option[string]{NewOption("field1 opt "+trigger, "field1 opt "+trigger)}
		}, &trigger)
	field2 := NewSelect[string]().
		TitleFunc(func() string {
			return "field2 title " + trigger
		}, &trigger).
		DescriptionFunc(func() string {
			return "field2 desc " + trigger
		}, &trigger).
		OptionsFunc(func() []Option[string] {
			return []Option[string]{NewOption("field2 opt "+trigger, "field2 opt "+trigger)}
		}, &trigger)
	field1.WithHeight(5)
	field2.WithHeight(5)
	f := NewForm(NewGroup(field1, field2)).WithHeight(10)

	doAllUpdates(f, f.Init())

	view := ansi.Strip(f.View())

	expectedStrings := []string{
		"field1 title initial",
		"field1 desc initial",
		"field1 opt initial",
		"field2 title initial",
		"field2 desc initial",
		"field2 opt initial",
	}
	for _, expected := range expectedStrings {
		if !strings.Contains(view, expected) {
			t.Log(pretty.Render(view))
			t.Error("Expected view to contain " + expected)
		}
	}

	if field1.GetValue() != "field1 opt initial" {
		t.Errorf("Expected field1 value to be field1 opt initial but was %s", field1.GetValue())
	}
	if field2.GetValue() != "field2 opt initial" {
		t.Errorf("Expected field2 value to be field2 opt initial but was %s", field2.GetValue())
	}

	trigger = "updated"
	_, cmd := f.Update(nil)
	doAllUpdates(f, cmd)
	view = ansi.Strip(f.View())

	expectedStrings = []string{
		"field1 title updated",
		"field1 desc updated",
		"field1 opt updated",
		"field2 title updated",
		"field2 desc updated",
		"field2 opt updated",
	}
	for _, expected := range expectedStrings {
		if !strings.Contains(view, expected) {
			t.Log(pretty.Render(view))
			t.Error("Expected view to contain " + expected)
		}
	}

	if field1.GetValue() != "field1 opt updated" {
		t.Errorf("Expected field1 value to be field1 opt updated but was %s", field1.GetValue())
	}
	if field2.GetValue() != "field2 opt updated" {
		t.Errorf("Expected field2 value to be field2 opt updated but was %s", field1.GetValue())
	}
}

func TestMultiSelect(t *testing.T) {
	field := NewMultiSelect[string]().
		Options(NewOptions(
			"Foo\nLine2",
			"Bar\nLine2",
			"Baz\nLine2",
			"Ban\nLine2",
		)...).
		Title("Which one?")
	f := NewForm(NewGroup(field)).
		WithHeight(5)
	f.Update(f.Init())

	view := ansi.Strip(f.View())

	if !strings.Contains(view, "Foo") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain Foo.")
	}

	if !strings.Contains(view, "Which one?") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain Which one?.")
	}

	if !strings.Contains(view, "> • Foo") {
		t.Log(pretty.Render(view))
		t.Error("Expected cursor to be on Foo.")
	}

	// Move selection cursor down
	m, _ := f.Update(keys('j'))
	view = ansi.Strip(m.View())

	if got, ok := field.Hovered(); !ok || got != "Bar\nLine2" {
		t.Log(pretty.Render(view))
		t.Error("Expected cursor to be on Bar.")
	}

	if strings.Contains(view, "> • Foo") {
		t.Log(pretty.Render(view))
		t.Error("Expected cursor to be on Bar.")
	}

	if !strings.Contains(view, "> • Bar") {
		t.Log(pretty.Render(view))
		t.Error("Expected cursor to be on Bar.")
	}

	// Toggle
	m, _ = f.Update(keys('x'))
	view = ansi.Strip(m.View())

	if !strings.Contains(view, "> ✓ Bar") {
		t.Log(pretty.Render(view))
		t.Error("Expected cursor to be on Bar.")
	}

	if !strings.Contains(view, "x toggle • ↑ up • ↓ down • / filter • enter submit") {
		t.Log(pretty.Render(view))
		t.Error("Expected field to contain help.")
	}

	// Submit
	f.Update(tea.KeyMsg{Type: tea.KeyEnter})

	value := field.GetValue()
	if value, ok := value.([]string); !ok {
		t.Error("Expected field value to a slice of string")
	} else {
		if len(value) != 1 {
			t.Error("Expected field value length to be 1")
		} else {
			if value[0] != "Bar\nLine2" {
				t.Error("Expected first field value to be Bar")
			}
		}
	}
}

func TestMultiSelectFiltering(t *testing.T) {
	tests := []struct {
		name      string
		filtering bool
	}{
		{"Filtering off", false},
		{"Filtering on", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			field := NewMultiSelect[string]().Options(NewOptions("Foo", "Bar", "Baz")...).Title("Which one?").Filterable(tc.filtering)
			f := NewForm(NewGroup(field))
			f.Update(f.Init())
			// Filter for values starting with a 'B' only.
			f.Update(keys('/'))
			f.Update(keys('B'))

			view := ansi.Strip(f.View())
			// When we're filtering, the list should change.
			if tc.filtering && strings.Contains(view, "Foo") {
				t.Log(pretty.Render(view))
				t.Error("Foo should not in filtered list.")
			}
			// When we're not filtering, the list shouldn't change.
			if !tc.filtering && !strings.Contains(view, "Foo") {
				t.Log(pretty.Render(view))
				t.Error("Expected list to contain Foo.")
			}
		})
	}
	t.Run("Remove filter option from help menu.", func(t *testing.T) {
		field := NewMultiSelect[string]().Options(NewOptions("Foo", "Bar", "Baz")...).Title("Which one?").Filterable(false)
		f := NewForm(NewGroup(field))
		f.Update(f.Init())
		view := ansi.Strip(f.View())
		if strings.Contains(view, "filter") {
			t.Log(pretty.Render(view))
			t.Error("Expected list to hide filtering in help menu.")
		}
	})
}

func TestSelectPageNavigation(t *testing.T) {
	opts := NewOptions(
		"Qux",
		"Quux",
		"Foo",
		"Bar",
		"Baz",
		"Corge",
		"Grault",
		"Garply",
		"Waldo",
		"Fred",
		"Plugh",
		"Xyzzy",
		"Thud",
		"Norf",
		"Blip",
		"Flob",
		"Zorp",
		"Smurf",
		"Bloop",
		"Ping",
	)

	reFirst := regexp.MustCompile(`>( •)? Qux`)
	reLast := regexp.MustCompile(`>( •)? Ping`)
	reHalfDown := regexp.MustCompile(`>( •)? Baz`)

	for name, field := range map[string]Field{
		"multiselect": NewMultiSelect[string]().Options(opts...).Title("Choose"),
		"select":      NewSelect[string]().Options(opts...).Title("Choose"),
	} {
		t.Run(name, func(t *testing.T) {
			f := NewForm(NewGroup(field)).WithHeight(10)
			f.Update(f.Init())

			view := ansi.Strip(f.View())
			if !reFirst.MatchString(view) {
				t.Log(pretty.Render(view))
				t.Errorf("Wrong item selected, should have matched %q (first item)", reFirst.String())
			}

			m, _ := f.Update(keys('G'))
			view = ansi.Strip(m.View())
			if !reLast.MatchString(view) {
				t.Log(pretty.Render(view))
				t.Errorf("Wrong item selected, should have matched %q (last item)", reLast.String())
			}

			m, _ = f.Update(keys('g'))
			view = ansi.Strip(m.View())
			if !reFirst.MatchString(view) {
				t.Log(pretty.Render(view))
				t.Errorf("Wrong item selected, should have matched %q (first item)", reFirst.String())
			}

			m, _ = f.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
			view = ansi.Strip(m.View())
			if !reHalfDown.MatchString(view) {
				t.Log(pretty.Render(view))
				t.Errorf("Wrong item selected, should have matched %q (half down item)", reHalfDown.String())
			}

			// sends multiple to verify it stays within boundaries
			f.Update(tea.KeyMsg{Type: tea.KeyCtrlU})
			f.Update(tea.KeyMsg{Type: tea.KeyCtrlU})
			m, _ = f.Update(tea.KeyMsg{Type: tea.KeyCtrlU})
			view = ansi.Strip(m.View())
			if !reFirst.MatchString(view) {
				t.Log(pretty.Render(view))
				t.Errorf("Wrong item selected, should have matched %q (first item)", reFirst.String())
			}

			// verify it stays within boundaries
			f.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
			f.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
			f.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
			f.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
			f.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
			m, _ = f.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
			view = ansi.Strip(m.View())
			if !reLast.MatchString(view) {
				t.Log(pretty.Render(view))
				t.Errorf("Wrong item selected, should have matched %q (last item)", reLast.String())
			}
		})
	}
}

func TestFile(t *testing.T) {
	field := NewFilePicker().Title("Which file?")
	cmd := field.Init()
	field.Update(cmd())

	view := ansi.Strip(field.View())

	if !strings.Contains(view, "No file selected") {
		t.Log(pretty.Render(view))
		t.Error("Expected file picker to show no file selected.")
	}

	if !strings.Contains(view, "Which file?") {
		t.Log(pretty.Render(view))
		t.Error("Expected file picker to show title.")
	}
}

func TestHideGroup(t *testing.T) {
	f := NewForm(
		NewGroup(NewNote().Description("Foo")).
			WithHide(true),
		NewGroup(NewNote().Description("Bar")),
		NewGroup(NewNote().Description("Baz")),
		NewGroup(NewNote().Description("Qux")).
			WithHideFunc(func() bool { return false }).
			WithHide(true),
	)

	f = batchUpdate(f, f.NextGroup()).(*Form)

	if v := f.View(); !strings.Contains(v, "Bar") {
		t.Log(pretty.Render(v))
		t.Error("expected Bar to be visible")
	}

	// should have no effect as previous group is hidden
	f.Update(prevGroup())

	if v := f.View(); !strings.Contains(v, "Bar") {
		t.Log(pretty.Render(v))
		t.Error("expected Bar to be visible")
	}

	f.Update(nextGroup())

	if v := f.View(); !strings.Contains(v, "Baz") {
		t.Log(pretty.Render(v))
		t.Error("expected Baz to be visible")
	}

	f.Update(nextGroup())

	if v := f.View(); strings.Contains(v, "Qux") {
		t.Log(pretty.Render(v))
		t.Error("expected Qux to be hidden")
	}

	if v := f.State; v != StateCompleted {
		t.Error("should have been completed")
	}
}

func TestHideGroupLastAndFirstGroupsNotHidden(t *testing.T) {
	f := NewForm(
		NewGroup(NewNote().Description("Bar")),
		NewGroup(NewNote().Description("Foo")).
			WithHide(true),
		NewGroup(NewNote().Description("Baz")),
	)

	f = batchUpdate(f, f.Init()).(*Form)

	if v := ansi.Strip(f.View()); !strings.Contains(v, "Bar") {
		t.Log(pretty.Render(v))
		t.Error("expected Bar to not be hidden")
	}

	// should have no effect as there isn't any
	f.Update(prevGroup())

	if v := f.View(); !strings.Contains(v, "Bar") {
		t.Log(pretty.Render(v))
		t.Error("expected Bar to not be hidden")
	}

	f.Update(nextGroup())

	if v := ansi.Strip(f.View()); !strings.Contains(v, "Baz") {
		t.Log(pretty.Render(v))
		t.Error("expected Baz to not be hidden")
	}

	// should submit the form
	f.Update(nextGroup())
	if v := f.State; v != StateCompleted {
		t.Error("should have been completed")
	}
}

func TestPrevGroup(t *testing.T) {
	f := NewForm(
		NewGroup(NewNote().Description("Bar")),
		NewGroup(NewNote().Description("Foo")),
		NewGroup(NewNote().Description("Baz")),
	)

	f = batchUpdate(f, f.Init()).(*Form)
	f.Update(nextGroup())
	f.Update(nextGroup())
	f.Update(prevGroup())
	f.Update(prevGroup())

	if v := ansi.Strip(f.View()); !strings.Contains(v, "Bar") {
		t.Log(pretty.Render(v))
		t.Error("expected Bar to not be hidden")
	}
}

func TestNote(t *testing.T) {
	field := NewNote().
		Title("Taco").
		Description("How may we take your order?").
		Next(true)
	f := NewForm(NewGroup(field))
	f.Update(f.Init())

	view := ansi.Strip(f.View())

	if !strings.Contains(view, "Taco") {
		t.Log(view)
		t.Error("Expected field to contain Taco title.")
	}

	if !strings.Contains(view, "order?") {
		t.Log(view)
		t.Error("Expected field to contain Taco description.")
	}

	if !strings.Contains(view, "Next") {
		t.Log(view)
		t.Error("Expected field to contain next button")
	}

	const expect = 7
	if h := lipgloss.Height(ansi.Strip(view)); h != expect {
		t.Log(view)
		t.Errorf("Expected field to have height %d, got %d", expect, h)
	}

	if !strings.Contains(view, "enter submit") {
		t.Log(view)
		t.Error("Expected field to contain help.")
	}
}

func TestDynamicHelp(t *testing.T) {
	f := NewForm(
		NewGroup(
			NewInput().Title("Dynamic Help"),
			NewInput().Title("Dynamic Help"),
			NewInput().Title("Dynamic Help"),
		),
	)
	f.Update(f.Init())

	view := ansi.Strip(f.View())

	if !strings.Contains(view, "Dynamic Help") {
		t.Log(pretty.Render(view))
		t.Fatal("Expected help to contain title.")
	}

	if strings.Contains(view, "shift+tab") || strings.Contains(view, "submit") {
		t.Log(pretty.Render(view))
		t.Error("Expected help not to contain shift+tab or submit.")
	}
}

func TestSkip(t *testing.T) {
	f := NewForm(
		NewGroup(
			NewInput().Title("First"),
			NewNote().Title("Skipped"),
			NewNote().Title("Skipped"),
			NewInput().Title("Second"),
		),
	).WithWidth(25)

	f = batchUpdate(f, f.Init()).(*Form)
	view := ansi.Strip(f.View())

	if !strings.Contains(view, "┃ First") {
		t.Log(pretty.Render(view))
		t.Error("Expected first field to be focused")
	}

	// next field should skip both of the notes and proceed to the last input.
	f.Update(NextField())
	view = ansi.Strip(f.View())

	if strings.Contains(view, "┃ First") {
		t.Log(pretty.Render(view))
		t.Error("Expected first field to be blurred")
	}

	if !strings.Contains(view, "┃ Second") {
		t.Log(pretty.Render(view))
		t.Error("Expected second field to be focused")
	}

	// previous field should skip both of the notes and focus the first input.
	f.Update(PrevField())
	view = ansi.Strip(f.View())

	if strings.Contains(view, "┃ Second") {
		t.Log(pretty.Render(view))
		t.Error("Expected second field to be blurred")
	}

	if !strings.Contains(view, "┃ First") {
		t.Log(pretty.Render(view))
		t.Error("Expected first field to be focused")
	}
}

func TestTimeout(t *testing.T) {
	// This test requires a real program, so make sure it doesn't interfere with our test runner.
	f := formProgram()

	// Test that the form times out after 1ms and returns a timeout error.
	err := f.WithTimeout(100 * time.Millisecond).Run()
	if err == nil || !errors.Is(err, ErrTimeout) {
		t.Errorf("expected timeout error, got %v", err)
	}
}

func TestAbort(t *testing.T) {
	// This test requires a real program, so make sure it doesn't interfere with our test runner.
	f := formProgram()

	// Test that the form aborts without throwing a timeout error when explicitly told to abort.
	ctx, cancel := context.WithCancel(context.Background())
	// Since the context is cancelled, the program should exit immediately.
	cancel()
	// Tell the form to abort.
	f.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	// Run the program.
	err := f.RunWithContext(ctx)
	if err == nil || !errors.Is(err, ErrUserAborted) {
		t.Errorf("expected user aborted error, got %v", err)
	}
}

const (
	title       = "A Title"
	description = "A Description"
)

var titleAndDescTests = map[string]struct {
	Empty       interface{ View() string }
	EmptyHeight int
	Title       interface{ View() string }
	Description interface{ View() string }
}{
	"Group": {
		NewGroup(NewInput()),
		1, // >
		NewGroup(NewInput()).Title(title),
		NewGroup(NewInput()).Description(description),
	},
	"Confirm": {
		NewConfirm(),
		1, // yes | no
		NewConfirm().Title(title),
		NewConfirm().Description(description),
	},
	"FilePicker": {
		NewFilePicker(),
		1, // "no file selected"
		NewFilePicker().Title(title),
		NewFilePicker().Description(description),
	},
	"Input": {
		NewInput(),
		1, // >
		NewInput().Title(title),
		NewInput().Description(description),
	},
	"Note": {
		NewNote(),
		1, // |
		NewNote().Title(title),
		NewNote().Description(description),
	},
	"Text": {
		NewText(),
		6, // textarea
		NewText().Title(title),
		NewText().Description(description),
	},
	"Select": {
		NewSelect[string](),
		1, // >
		NewSelect[string]().Title(title),
		NewSelect[string]().Description(description),
	},
	"MultiSelect": {
		NewMultiSelect[string](),
		1, // >
		NewMultiSelect[string]().Title(title),
		NewMultiSelect[string]().Description(description),
	},
}

func TestNoTitleOrDescription(t *testing.T) {
	for name, tt := range titleAndDescTests {
		t.Run(name, func(t *testing.T) {
			view := tt.Empty.View()
			got := lipgloss.Height(ansi.Strip(view))
			want := tt.EmptyHeight
			if got != want {
				t.Log(pretty.Render(view))
				t.Fatalf("got != want; height should be %d, got %d", want, got)
			}
		})
	}
}

func TestTitleRowRender(t *testing.T) {
	for name, tt := range titleAndDescTests {
		t.Run(name, func(t *testing.T) {
			view := tt.Title.View()
			if !strings.Contains(view, title) {
				t.Log(pretty.Render(view))
				t.Error("Expected title to be visible")
			}
		})
	}
}

func TestDescriptionRowRender(t *testing.T) {
	for name, tt := range titleAndDescTests {
		t.Run(name, func(t *testing.T) {
			view := tt.Description.View()
			if !strings.Contains(view, description) {
				t.Log(pretty.Render(view))
				t.Error("Expected description to be visible")
			}
		})
	}
}

func TestGetFocusedField(t *testing.T) {
	f := NewForm(
		NewGroup(
			NewInput().Title("First").Key("First"),
			NewInput().Title("Second").Key("Second"),
			NewInput().Title("Third").Key("Third"),
		),
	).WithWidth(25)
	f = batchUpdate(f, f.Init()).(*Form)

	f.NextField()
	field := f.GetFocusedField()

	if field.GetKey() != "Second" {
		t.Error("Expected Second field to be focused but was '" + field.GetKey() + "'")
	}
}

// formProgram returns a new Form with a nil input and output, so it can be used as a test program.
func formProgram() *Form {
	return NewForm(NewGroup(NewInput().Title("Foo"))).
		WithInput(nil).
		WithOutput(io.Discard).
		WithAccessible(false)
}

func batchUpdate(m tea.Model, cmd tea.Cmd) tea.Model {
	if cmd == nil {
		return m
	}
	msg := cmd()
	m, cmd = m.Update(msg)
	if cmd == nil {
		return m
	}
	msg = cmd()
	m, _ = m.Update(msg)
	return m
}

func keys(runes ...rune) tea.KeyMsg {
	return tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: runes,
	}
}

func TestAccessibleForm(t *testing.T) {
	var out bytes.Buffer

	f := NewForm(
		NewGroup(
			NewInput().Title("Hello:"),
		),
	).
		WithAccessible(true).
		WithOutput(&out).
		WithInput(strings.NewReader("carlos\n"))

	if err := f.Run(); err != nil {
		t.Error(err)
	}

	if !strings.Contains(out.String(), "Hello: ") {
		t.Error("invalid output:\n", out.String())
	}
}

func TestAccessibleFields(t *testing.T) {
	for name, test := range map[string]struct {
		Field       Field
		FieldFn     func() Field
		Input       string
		CheckOutput func(tb testing.TB, output string)
		CheckValue  func(tb testing.TB, value any)
	}{
		"input": {
			Field: NewInput(),
			Input: "Hello",
			CheckOutput: func(tb testing.TB, output string) {
				tb.Helper()
				requireContains(tb, output, "Input: ")
			},
			CheckValue: func(tb testing.TB, value any) {
				tb.Helper()
				requireEqual(tb, "Hello", value.(string))
			},
		},
		"input with charlimit": {
			Field: NewInput().CharLimit(2),
			Input: "Hello",
			CheckOutput: func(tb testing.TB, output string) {
				tb.Helper()
				requireContains(tb, output, "Input cannot exceed 2 characters")
			},
		},
		"input with default": {
			FieldFn: func() Field {
				v := "hi"
				return NewInput().Value(&v)
			},
			Input: "\n",
			CheckOutput: func(tb testing.TB, output string) {
				tb.Helper()
				requireContains(tb, output, "Input: ")
			},
			CheckValue: func(tb testing.TB, value any) {
				tb.Helper()
				requireEqual(tb, "hi", value.(string))
			},
		},
		"confirm": {
			Field: NewConfirm(),
			Input: "Y",
			CheckOutput: func(tb testing.TB, output string) {
				tb.Helper()
				requireContains(tb, output, "Choose [y/N] ")
			},
			CheckValue: func(tb testing.TB, value any) {
				tb.Helper()
				requireEqual(tb, true, value.(bool))
			},
		},
		"confirm with default": {
			FieldFn: func() Field {
				v := true
				return NewConfirm().Value(&v)
			},
			CheckOutput: func(tb testing.TB, output string) {
				tb.Helper()
				requireContains(tb, output, "Choose [Y/n] ")
			},
			CheckValue: func(tb testing.TB, value any) {
				tb.Helper()
				requireEqual(tb, true, value.(bool))
			},
		},
		"confirm with default choose": {
			FieldFn: func() Field {
				v := true
				return NewConfirm().Value(&v)
			},
			Input: "n",
			CheckOutput: func(tb testing.TB, output string) {
				tb.Helper()
				requireContains(tb, output, "Y/n")
			},
			CheckValue: func(tb testing.TB, value any) {
				tb.Helper()
				requireEqual(tb, false, value.(bool))
			},
		},
		"filepicker": {
			Field: NewFilePicker(),
			Input: "huh_test.go",
			CheckOutput: func(tb testing.TB, output string) {
				tb.Helper()
				requireContains(tb, output, "Choose a file: ")
			},
			CheckValue: func(tb testing.TB, value any) {
				tb.Helper()
				requireEqual(tb, "huh_test.go", value.(string))
			},
		},
		"filepicker with default": {
			FieldFn: func() Field {
				v := "huh_test.go"
				return NewFilePicker().Value(&v)
			},
			Input: "\n",
			CheckOutput: func(tb testing.TB, output string) {
				tb.Helper()
				requireContains(tb, output, "Choose a file: ")
			},
			CheckValue: func(tb testing.TB, value any) {
				tb.Helper()
				requireEqual(tb, "huh_test.go", value.(string))
			},
		},
		"multiselect": {
			Field: NewMultiSelect[string]().Options(NewOptions("a", "b")...),
			Input: "2",
			CheckOutput: func(tb testing.TB, output string) {
				tb.Helper()
				requireContains(tb, output, "2. ✓ b")
			},
			CheckValue: func(tb testing.TB, value any) {
				tb.Helper()
				got := value.([]string)
				requireEqual(tb, 1, len(got))
				requireEqual(tb, "b", got[0])
			},
		},
		"multiselect default value": {
			FieldFn: func() Field {
				v := []string{"b", "c"}
				return NewMultiSelect[string]().Options(NewOptions("a", "b", "c", "d")...).Value(&v)
			},
			Input: "\n",
			CheckOutput: func(tb testing.TB, output string) {
				tb.Helper()
				requireContains(tb, output, "2. ✓ b")
				requireContains(tb, output, "3. ✓ c")
			},
		},
		"select": {
			Field: NewSelect[string]().Options(NewOptions("a", "b")...),
			Input: "2",
			CheckOutput: func(tb testing.TB, output string) {
				tb.Helper()
				requireContains(tb, output, "Select: ")
			},
			CheckValue: func(tb testing.TB, value any) {
				tb.Helper()
				requireEqual(tb, "b", value.(string))
			},
		},
		"select default value": {
			FieldFn: func() Field {
				v := "c"
				return NewSelect[string]().
					Options(NewOptions("a", "b", "c", "d")...).
					Value(&v)
			},
			Input: "\n",
			CheckOutput: func(tb testing.TB, output string) {
				tb.Helper()
				requireContains(tb, output, "Select: ")
				requireContains(tb, output, "Enter a number between 1 and 4")
			},
			CheckValue: func(tb testing.TB, value any) {
				tb.Helper()
				requireEqual(tb, "c", value.(string))
			},
		},
		"select no input": {
			Field: NewSelect[string]().Options(NewOptions("a", "b")...),
			Input: "\n2\n",
			CheckOutput: func(tb testing.TB, output string) {
				tb.Helper()
				requireContains(tb, output, "Select: ")
				requireContains(tb, output, "Enter a number between 1 and 2")
				requireContains(tb, output, "Invalid: must be a number between 1 and 2")
			},
			CheckValue: func(tb testing.TB, value any) {
				tb.Helper()
				requireEqual(tb, "b", value.(string))
			},
		},
		"select single option": {
			Field: NewSelect[string]().Options(NewOptions("a")...),
			Input: "\n1\n",
			CheckOutput: func(tb testing.TB, output string) {
				tb.Helper()
				requireContains(tb, output, "Select: ")
				requireContains(tb, output, "There is only one option available; enter the number 1:")
				requireContains(tb, output, "Invalid: must be 1")
			},
			CheckValue: func(tb testing.TB, value any) {
				tb.Helper()
				requireEqual(tb, "a", value.(string))
			},
		},
		"note": {
			Field: NewNote().Title("Hi").Description("there"),
			CheckOutput: func(tb testing.TB, output string) {
				tb.Helper()
				requireContains(tb, output, "Hi")
				requireContains(tb, output, "there")
			},
		},
		"text": {
			Field: NewText().Title("Text: "),
			Input: "hello world",
			CheckOutput: func(tb testing.TB, output string) {
				tb.Helper()
				requireContains(tb, output, "Text: ")
			},
			CheckValue: func(tb testing.TB, value any) {
				tb.Helper()
				requireEqual(tb, "hello world", value.(string))
			},
		},
		"text with limit": {
			Field: NewText().CharLimit(2).Title("Text"),
			Input: "hello world",
			CheckOutput: func(tb testing.TB, output string) {
				tb.Helper()
				requireContains(tb, output, "Input cannot exceed 2 characters")
			},
		},
		"text default value": {
			FieldFn: func() Field {
				v := "test"
				return NewText().Title("Text:").Value(&v)
			},
			Input: "\n",
			CheckOutput: func(tb testing.TB, output string) {
				tb.Helper()
				requireContains(tb, output, "Text: ")
			},
			CheckValue: func(tb testing.TB, value any) {
				tb.Helper()
				requireEqual(tb, "test", value.(string))
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			field := test.Field
			if test.FieldFn != nil {
				field = test.FieldFn()
			}

			var out bytes.Buffer
			if err := field.RunAccessible(
				&out,
				strings.NewReader(test.Input),
			); err != nil {
				t.Error(err)
			}
			if test.CheckOutput != nil {
				test.CheckOutput(t, out.String())
			}
			if test.CheckValue != nil {
				test.CheckValue(t, field.GetValue())
			}
		})
	}
}

func TestInputPasswordAccessible(t *testing.T) {
	t.Run("not a tty", func(t *testing.T) {
		var out bytes.Buffer
		if err := NewInput().
			EchoMode(EchoModeNone).
			RunAccessible(&out, bytes.NewReader(nil)); err == nil {
			t.Error("expected it to error")
		}
		if err := NewInput().
			EchoMode(EchoModePassword).
			RunAccessible(&out, bytes.NewReader(nil)); err == nil {
			t.Error("expected it to error")
		}
	})

	t.Run("is a tty", func(t *testing.T) {
		var out bytes.Buffer
		pty, err := xpty.NewPty(50, 30)
		if err != nil {
			t.Skipf("could not open pty: %v", err)
		}
		upty, ok := pty.(*xpty.UnixPty)
		if !ok {
			t.Skipf("test only works on unix")
		}

		input := NewInput().EchoMode(EchoModePassword)

		errs := make(chan error, 1)
		go func() {
			errs <- input.RunAccessible(&out, upty.Slave())
		}()

		upty.Master().Write([]byte("a password\n"))

		if err := <-errs; err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		requireContains(t, out.String(), "Password: ")
		requireEqual(t, "a password", input.GetValue().(string))
	})
}

func requireEqual[T comparable](tb testing.TB, a, b T) {
	tb.Helper()
	if a != b {
		tb.Fatalf("expected %v to be equal to %v", a, b)
	}
}

func requireContains(tb testing.TB, s, subtr string) {
	tb.Helper()
	if !strings.Contains(s, subtr) {
		tb.Fatalf("%q does not contain %q", s, subtr)
	}
}
