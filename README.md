# Huh?

A simple and powerful library for building interactive forms in the terminal. Powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea).

<img alt="Running a taco form" width="600" src="./examples/taco/taco.gif">

The above example is running from a single Go program ([source](./examples/taco/main.go)).

## Tutorial

`huh?` provides a straightforward API to build forms and prompt users for input.

For this tutorial, we're building a Taco order form. Lets start by importing the
dependencies that we'll need.

```go
package main

import (
  "log"

  "github.com/charmbracelet/huh"
)
```

`huh` allows you to define a form with multiple groups to separate field forms
into pages. We will set up a form with three groups for the customer to fill
out.

```go
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
```

Finally, we can run the form:

```go
err := form.Run()
if err != nil {
  log.Fatal(err)
}
```

## Field Reference
* [`Input`](#input)
* [`Text`](#text)
* [`Select`](#select)
* [`MultiSelect`](#multi-select)
* [`Confirm`](#confirm)

### Input

```go
huh.NewInput().
	Title("What's for lunch?").
	Prompt("?").
	Validate(isFood).
	Value(&lunch)
```

### Text

```go
huh.NewText().
	Title("Tell me a story.").
	Validate(checkForPlagiarism).
	Value(&story)
```

### Select

```go
huh.NewSelect[string]().
	Title("Pick a country.").
	Options(
		huh.NewOption("United States", "US"),
		huh.NewOption("Germany", "DE"),
		huh.NewOption("Brazil", "BR"),
		huh.NewOption("Canada", "CA"),
	).
	Value(&country)
```

### Multiple Select

```go
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
```

### Confirm

```go
huh.NewConfirm().
  Title("Toppings").
  Affirmative("Yes!").
  Negative("No.").
  Value(&confirm)
```

## Feedback

We'd love to hear your thoughts on this project. Feel free to drop us a note!

* [Twitter](https://twitter.com/charmcli)
* [The Fediverse](https://mastodon.social/@charmcli)
* [Discord](https://charm.sh/chat)

## Acknowledgments

`huh` is inspired by the wonderful [Survey][survey] library by Alec Aivazis.

[survey]: https://github.com/AlecAivazis/survey

## License

[MIT](https://github.com/charmbracelet/bubbletea/raw/master/LICENSE)

***

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400"></a>

Charm热爱开源 • Charm loves open source • نحنُ نحب المصادر المفتوحة
