package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	xstrings "github.com/charmbracelet/x/exp/strings"
)

type Spice int

const (
	Mild Spice = iota + 1
	Medium
	Hot
)

func (s Spice) String() string {
	switch s {
	case Mild:
		return "mild "
	case Medium:
		return "medium-spicy "
	case Hot:
		return "hot "
	default:
		return ""
	}
}

type Order struct {
	Taco         Taco
	Name         string
	Instructions string
	Discount     bool
}

type Taco struct {
	Shell    string
	Spice    Spice
	Base     string
	Toppings []string
}

func main() {
	var taco Taco
	var order = Order{Taco: taco, Instructions: "Dressing on the side"}

	// Should we run in accessible mode?
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	form := huh.NewForm(
		huh.NewGroup(huh.NewNote().
			Title("Taco Charm").
			Description("Welcome to _Taco Charm™_.\n\nHow may we take your order?").
			Next(true)),

		// What's a taco without a shell?
		// We'll need to know what filling to put inside too.
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(huh.NewOptions("Soft", "Hard")...).
				Title("Shell?").
				Description("Our tortillas are made fresh in-house, every day.").
				Validate(func(t string) error {
					if t == "Hard" {
						return fmt.Errorf("we're out of hard shells, sorry")
					}
					return nil
				}).
				Value(&order.Taco.Shell),

			huh.NewSelect[string]().
				Options(huh.NewOptions("Chicken", "Beef", "Fish", "Beans")...).
				Value(&order.Taco.Base).
				Title("Base"),
		),

		// Prompt for toppings and special instructions.
		// The customer can ask for up to 4 toppings.
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Toppings").
				Description("Choose up to 4.").
				Options(
					huh.NewOption("Lettuce", "lettuce").Selected(true),
					huh.NewOption("Tomatoes", "tomatoes").Selected(true),
					huh.NewOption("Corn", "corn"),
					huh.NewOption("Salsa", "salsa"),
					huh.NewOption("Sour Cream", "sour cream"),
					huh.NewOption("Cheese", "cheese"),
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

			huh.NewSelect[Spice]().
				Title("Spice Level").
				Options(
					huh.NewOption("Mild", Mild),
					huh.NewOption("Medium", Medium),
					huh.NewOption("Hot", Hot).Selected(true),
				).
				Value(&order.Taco.Spice),
		),

		// Gather final details for the order.
		huh.NewGroup(
			huh.NewInput().
				Value(&order.Name).
				Title("What's your name?").
				Placeholder("Margaret Thatcher").
				Description("For when your order is ready."),

			huh.NewText().
				Value(&order.Instructions).
				Placeholder("Just put it in the mailbox please").
				Title("Special Instructions").
				Description("Anything we should know?").
				CharLimit(400),

			huh.NewConfirm().
				Title("Would you like 15% off?").
				Value(&order.Discount).
				Affirmative("Yes!").
				Negative("No."),
		),
	).WithAccessible(accessible)

	err := form.Run()

	if err != nil {
		log.Fatal(err)
	}

	prepareTaco := func() {
		time.Sleep(2 * time.Second)
	}

	_ = spinner.New().Title("Preparing your taco").Static(accessible).Action(prepareTaco).Run()

	// Print order summary.
	{
		var sb strings.Builder
		keyword := func(s string) string {
			return strings.ToLower(lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s))
		}
		fmt.Fprintf(&sb,
			"%s\n\nOne %s%s shell taco filled with %s and topped with %s.",
			lipgloss.NewStyle().Bold(true).Render("TACO RECEIPT"),
			keyword(order.Taco.Spice.String()),
			keyword(order.Taco.Shell),
			keyword(order.Taco.Base),
			keyword(xstrings.EnglishJoin(order.Taco.Toppings, true)),
		)

		name := order.Name
		if name != "" {
			name = ", " + name
		}
		fmt.Fprintf(&sb, "\n\nThanks for your order%s!", name)

		if order.Discount {
			fmt.Fprint(&sb, "\n\nEnjoy 15% off.")
		}

		fmt.Println(
			lipgloss.NewStyle().
				Width(40).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63")).
				Padding(1, 2).
				Render(sb.String()),
		)
	}
}
