package main

import (
	"log"

	"github.com/charmbracelet/huh"
)

// TODO: ensure input is not plagiarized.
func checkForPlagiarism(s string) error { return nil }

// TODO: ensure input is food.
func isFood(s string) error { return nil }

// TODO: ensure input is a valid name.
func validateName(s string) error { return nil }

func main() {
	var (
		lunch    string
		story    string
		country  string
		toppings []string
		discount bool
	)

	// `Input`s are single line text fields.
	huh.NewInput().
		Title("What's for lunch?").
		Prompt("?").
		Validate(isFood).
		Value(&lunch)

	// `Text`s are multi-line text fields.
	huh.NewText().
		Title("Tell me a story.").
		Validate(checkForPlagiarism).
		Value(&story)

	// `Select`s are multiple choice questions.
	huh.NewSelect[string]().
		Title("Pick a country.").
		Options(
			huh.NewOption("United States", "US"),
			huh.NewOption("Germany", "DE"),
			huh.NewOption("Brazil", "BR"),
			huh.NewOption("Canada", "CA"),
		).
		Value(&country)

	// `MultiSelect`s allow multiple selections from a list of options.
	huh.NewMultiSelect[string]().
		Options(
			huh.NewOption("Cheese", "cheese").Selected(true),
			huh.NewOption("Lettuce", "lettuce").Selected(true),
			huh.NewOption("Corn", "corn"),
			huh.NewOption("Salsa", "salsa"),
			huh.NewOption("Sour Cream", "sour cream"),
			huh.NewOption("Tomatoes", "tomatoes"),
		).
		Title("Toppings").
		Limit(4).
		Value(&toppings)

	// `Confirm`s are a confirmation prompt.
	huh.NewConfirm().
		Title("Want a discount?").
		Affirmative("Yes!").
		Negative("No.").
		Value(&discount)

	// Form
	var (
		shell        string
		base         string
		name         string
		instructions string
	)

	form := huh.NewForm(
		// What's a taco without a shell?
		// We'll need to know what filling to put inside too.
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(
					huh.NewOption("Soft", "soft"),
					huh.NewOption("Hard", "hard"),
				).
				Title("Shell?").
				Value(&shell),

			huh.NewSelect[string]().
				Options(
					huh.NewOption("Chicken", "chicken"),
					huh.NewOption("Beef", "beef"),
					huh.NewOption("Fish", "fish"),
					huh.NewOption("Beans", "beans"),
				).
				Title("Base").
				Value(&base),
		),

		// Prompt for toppings and special instructions.
		// The customer can ask for up to 4 toppings.
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Options(
					huh.NewOption("Tomatoes", "tomatoes").Selected(true),
					huh.NewOption("Lettuce", "lettuce").Selected(true),
					huh.NewOption("Salsa", "salsa"),
					huh.NewOption("Cheese", "cheese"),
					huh.NewOption("Sour Cream", "sour cream"),
					huh.NewOption("Corn", "corn"),
				).
				Title("Toppings").
				Limit(4).
				Value(&toppings),
		),

		// Gather final details for the order.
		huh.NewGroup(
			huh.NewInput().
				Title("What's your name?").
				Value(&name).
				Validate(validateName),

			huh.NewText().
				Title("Special Instructions").
				Value(&instructions).
				CharLimit(400),

			huh.NewConfirm().
				Title("Would you like 15% off").
				Value(&discount),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
}
