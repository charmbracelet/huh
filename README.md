# Huh?

A simple and powerful library for building interactive forms in the terminal. Powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea).

<img alt="Running a taco form" width="600" src="./examples/taco/taco.gif">

The above example is running from a single Go program ([source](./examples/taco/main.go)).

## Tutorial

`Huh?` provides a straightforward API to build forms and prompt users for input.

Let's build a simple form to take a order from a Taco shop.

Start by `import`ing `huh` into your Go program.

```go
package main

import  "github.com/charmbracelet/huh"
```

Create fields for the order:

```go
var (
  base     string
  name     string
  shell    string
  toppings []string
)

shellField := huh.NewSelect[string]().
  Options(huh.NewOptions("Hard", "Soft")...).
  Title("Shell?").
  Value(&shell)

baseField := huh.NewSelect[string]().
  Options(huh.NewOptions("Chicken", "Beef", "Fish", "Beans")...).
  Title("Base?").
  Value(&base)

toppingsField := huh.NewMultiSelect[string]().
  Options(huh.NewOptions("Tomatoes", "Lettuce", "Cheese", "Salsa", "Sour Cream", "Corn")...).
  Title("Toppings").
  Limit(4).
  Value(&toppings)

nameField := huh.NewInput().
  Title("Name").
  Value(&name)
```

Create the form and group the fields:

```go
form := huh.NewForm(
  huh.NewGroup(shellField, baseField),
  huh.NewGroup(toppingsField),
  huh.NewGroup(nameField),
)
```

Run the form:

```go
err := form.Run()
if err != nil {
  log.Fatal(err)
}
```

Use the values:

```go
fmt.Println("Shell: ", shell)
fmt.Println("Base: ", base)
fmt.Println("Toppings: ", strings.Join(toppings, ", "))
fmt.Println("Name", name)
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
huh.NewSelect[string]().
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
huh.NewSelect[string]().
  Title("Pick a country.").
  Options(huh.NewOptions("United States", "Germany", "Brazil", "Canada")...).
  Cursor("→")
```

## Multiple Select

`MultiSelect`s are multiple choice questions but allow multiple selections.

```go
huh.NewMultiSelect[string]().
  Options(
    huh.NewOption("Cheese", "cheese").Selected(true),
    huh.NewOption("Lettuce", "lettuce").Selected(true),
    huh.NewOption("Tomatoes", "tomatoes"),
    huh.NewOption("Corn", "corn"),
    huh.NewOption("Salsa", "salsa"),
    huh.NewOption("Sour Cream", "sour cream"),
  ).
  Title("Toppings").
  Limit(4),
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
