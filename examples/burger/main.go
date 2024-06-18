package main

import (
	"errors"
	"fmt"
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
		return "Mild "
	case Medium:
		return "Medium-Spicy "
	case Hot:
		return "Spicy-Hot "
	default:
		return ""
	}
}

type Order struct {
	Burger       Burger
	Side         string
	Name         string
	Instructions string
	Discount     bool
}

type Burger struct {
	Type     string
	Toppings []string
	Spice    Spice
}

func main() {
	var burger Burger
	var order = Order{Burger: burger}

	// Should we run in accessible mode?
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	form := huh.NewForm(
		huh.NewGroup(huh.NewNote().
			Title("Charmburger").
			Description("Welcome to _Charmburger™_.\n\nHow may we take your order?\n\n").
			Next(true).
			NextLabel("Next"),
		),

		// Choose a burger.
		// We'll need to know what topping to add too.
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(huh.NewOptions("Charmburger Classic", "Chickwich", "Fishburger", "Charmpossible™ Burger")...).
				Title("Choose your burger").
				Description("At Charm we truly have a burger for everyone.").
				Validate(func(t string) error {
					if t == "Fishburger" {
						return fmt.Errorf("no fish today, sorry")
					}
					return nil
				}).
				Value(&order.Burger.Type),

			huh.NewMultiSelect[string]().
				Title("Toppings").
				Description("Choose up to 4.").
				Options(
					huh.NewOption("Lettuce", "Lettuce").Selected(true),
					huh.NewOption("Tomatoes", "Tomatoes").Selected(true),
					huh.NewOption("Charm Sauce", "Charm Sauce"),
					huh.NewOption("Jalapeños", "Jalapeños"),
					huh.NewOption("Cheese", "Cheese"),
					huh.NewOption("Vegan Cheese", "Vegan Cheese"),
					huh.NewOption("Nutella", "Nutella"),
				).
				Validate(func(t []string) error {
					if len(t) <= 0 {
						return fmt.Errorf("at least one topping is required")
					}
					return nil
				}).
				Value(&order.Burger.Toppings).
				Filterable(true).
				Limit(4),
		),

		// Prompt for toppings and special instructions.
		// The customer can ask for up to 4 toppings.
		huh.NewGroup(
			huh.NewSelect[Spice]().
				Title("Spice level").
				Options(
					huh.NewOption("Mild", Mild).Selected(true),
					huh.NewOption("Medium", Medium),
					huh.NewOption("Hot", Hot),
				).
				Value(&order.Burger.Spice),

			huh.NewSelect[string]().
				Options(huh.NewOptions("Fries", "Disco Fries", "R&B Fries", "Carrots")...).
				Value(&order.Side).
				Title("Sides").
				Description("You get one free side with this order."),
		),

		// Gather final details for the order.
		huh.NewGroup(
			huh.NewInput().
				Value(&order.Name).
				Title("What's your name?").
				Placeholder("Margaret Thatcher").
				Validate(func(s string) error {
					if s == "Frank" {
						return errors.New("no franks, sorry")
					}
					return nil
				}).
				Description("For when your order is ready."),

			huh.NewText().
				Value(&order.Instructions).
				Placeholder("Just put it in the mailbox please").
				Title("Special Instructions").
				Description("Anything we should know?").
				CharLimit(400).
				Lines(5),

			huh.NewConfirm().
				Title("Would you like 15% off?").
				Value(&order.Discount).
				Affirmative("Yes!").
				Negative("No."),
		),
	).WithAccessible(accessible)

	err := form.Run()

	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	prepareBurger := func() {
		time.Sleep(2 * time.Second)
	}

	_ = spinner.New().Title("Preparing your burger...").Accessible(accessible).Action(prepareBurger).Run()

	// Print order summary.
	{
		var sb strings.Builder
		keyword := func(s string) string {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
		}
		fmt.Fprintf(&sb,
			"%s\n\nOne %s%s, topped with %s with %s on the side.",
			lipgloss.NewStyle().Bold(true).Render("BURGER RECEIPT"),
			keyword(order.Burger.Spice.String()),
			keyword(order.Burger.Type),
			keyword(xstrings.EnglishJoin(order.Burger.Toppings, true)),
			keyword(order.Side),
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
