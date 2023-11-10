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

<img alt="Input field" width="600" src="./examples/readme/input.gif">

```go
huh.NewInput().
    Title("What's for lunch?").
    Prompt("?").
    Validate(isFood).
    Value(&lunch)
```

### Text

<img alt="Text field" width="600" src="./examples/readme/text.gif">

```go
huh.NewText().
    Title("Tell me a story.").
    Validate(checkForPlagiarism).
    Value(&story)
```

### Select

<img alt="Select field" width="600" src="./examples/readme/select.gif">

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

<img alt="Multiselect field" width="600" src="./examples/readme/multiselect.gif">

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

<img alt="Confirm field" width="600" src="./examples/readme/confirm.gif">

```go
huh.NewConfirm().
    Title("You sure?").
    Affirmative("Yes!").
    Negative("No.").
    Value(&confirm)
```

## Themes

Huh forms can be customized through themes. You can supply your own custom theme
or use the predefined themes.

There are currently four predefined themes:

* `Charm`
* `Dracula`
* `Base 16`
* `Default`

| Theme         | Image |
|--------------|:-----:|
| Charm | <img alt="Charm-themed form" width="400" src="./examples/theme/charm-theme.png"> |
| Dracula | <img alt="Dracula-themed form" width="400" src="./examples/theme/dracula-theme.png"> |
| Base 16 | <img alt="Base 16-themed form" width="400" src="./examples/theme/basesixteen-theme.png"> |
| Default | <img alt="Default-themed form" width="400" src="./examples/theme/default-theme.png"> |

## Spinner

Huh additionally provides a `spinner` subpackage for displaying spinners while
performing actions. It's useful to complete an action after your user completes
a form.

<img alt="Spinner while making a taco" width="600" src="./spinner/examples/loading/spinner.gif">

To get started, create a new spinner, set a title, set an action, and run the
spinner:

```go
makeTaco := func() {
    //...
}

err := spinner.New().
    Title("Making your taco...").
    Action(makeTaco).
    Run()

fmt.Println("Order up!")
```

Alternatively, you can also use `Context`s. The spinner will stop once the
context is cancelled.

```go
makeTaco := func() {
    // ...
}

ctx, _ := context.WithTimeout(context.Background(), time.Second)

go makeTaco()

err := spinner.New().
    Title("Making your taco...").
    Context(ctx).
    Run()

fmt.Println("Order up!")
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
