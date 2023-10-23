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
  "fmt"
  "log"

  "github.com/charmbracelet/huh"
)

func main() {
  form := huh.NewForm(
    // What's a taco without a shell?
    // We'll need to know what filling to put inside too.
    huh.Group(
      huh.Select().
        Title("Shell?").
        Options("Hard", "Soft"),

      huh.Select().
        Title("Base").
        Options("Chicken", "Beef", "Fish", "Beans"),
    ),

    // Prompt for toppings and special instructions.
    // The customer can ask for up to 4 toppings.
    huh.Group(
      huh.MultiSelect().
        Title("Toppings").
        Options("Lettuce", "Tomatoes", "Corn", "Salsa", "Sour Cream", "Cheese").
        Filterable(true).
        Limit(4),

      huh.Text().
        Title("Special Instructions").
        CharLimit(400),
      ),

    // Gather final details for the order.
    huh.Group(
      huh.Input().
        Key("name").
        Title("What's your name?").
        Validate(huh.ValidateLength(0, 20)),

      huh.Confirm().
        Key("discount").
        Title("Would you like 15% off"),
      ),
  )

  r, err := form.Run()
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("A %s shell filled with %s and %s, topped with %s.",
    r["Shell?"], r["Base"], r["Toppings"])

  fmt.Println("That will be $%.2f. Thanks for your order, %s!",
    calculatePrice(r), r["name"])
}
```

## Input

`Input`s are single line text fields.

```go
huh.Input().
  Title("What's for lunch?").
  Validate(huh.ValidateLength(0, 20)).
  Prompt("?")
```

## Text

`Text`s are multi-line text fields.

```go
huh.Text().
  Title("Tell me a story.").
  Validate(huh.ValidateLength(100, 400)).
  Prompt(">").
  Editor(true) // open in $EDITOR
```

## Select

`Select`s are multiple choice questions.

```go
huh.Select().
  Title("Pick a country.").
  Option("United States").
  Option("Germany").
  Option("Brazil").
  Option("Canada").
  Cursor("→")
```

Alternatively,

```go
huh.Select().
  Title("Pick a country.").
  Options("United States", "Germany", "Brazil", "Canada").
  Cursor("→")
```

## Multiple Select

`MultiSelect`s are multiple choice questions but allow multiple selections.

```go
huh.MultiSelect().
  Title("Toppings.").
  Option("Lettuce").
  Option("Tomatoes").
  Option("Cheese").
  Option("Corn").
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
