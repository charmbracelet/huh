# Huh?

A simple and powerful library for building interactive forms in the terminal. Powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Tutorial

`Huh?` provides a straightforward interface to prompt users for input.

There are several `Question` types available to use:
* [`Input`](#input)
* [`Text`](#text)
* [`Select`](#select)
* [`MultiSelect`](#multiple-select)

```go
package main

import (
  "log"

  "github.com/charmbracelet/huh"
)

func main() {
  form := huh.NewForm(
    // What's a taco without a shell?
    // We'll need to know what filling to put inside too.
    huh.NewGroup(
      huh.NewSelect("Hard", "Soft").
        Title("Shell?"),

      huh.NewSelect("Chicken", "Beef", "Fish", "Beans").
        Title("Base"),
    ),

    // Prompt for toppings and special instructions.
    // The customer can ask for up to 4 toppings.
    huh.NewGroup(
      huh.NewMultiSelect("Lettuce", "Tomatoes", "Corn", "Salsa", "Sour Cream", "Cheese").
        Title("Toppings").
        Limit(4),

      huh.NewText().
        Title("Special Instructions").
        CharLimit(400),
    ),

    // Gather final details for the order.
    huh.NewGroup(
      huh.NewInput().
        Title("What's your name?").
        Validate(validateName),

      huh.NewConfirm().
        Title("Would you like 15% off"),
    ),
  )

  err := form.Run()
  if err != nil {
    log.Fatal(err)
  }
}
```

## Input

`Input`s are single line text fields.

```go
huh.NewInput().
  Title("What's for lunch?").
  Validate(validateLength).
  Prompt("?")
```

## Text

`Text`s are multi-line text fields.

```go
huh.NewText().
  Title("Tell me a story.").
  Validate(validateLength).
  Prompt(">").
  Editor(true) // open in $EDITOR
```

## Select

`Select`s are multiple choice questions.

```go
huh.NewSelect[Country]().
  Title("Pick a country.").
  Options(
    huh.NewOption("United States", "US"),
    huh.NewOption("Germany", "DE"),
    huh.NewOption("Brazil", "BR"),
    huh.NewOption("Canada", "CA"),
  ).
  Cursor("→")
```

Alternatively,

```go
huh.NewSelect("United States", "Germany", "Brazil", "Canada").
  Title("Pick a country.").
  Cursor("→")
```

## Multiple Select

`MultiSelect`s are multiple choice questions but allow multiple selections.

```go
huh.NewMultiSelect().
  Title("Toppings.").
  Limit(4)
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
