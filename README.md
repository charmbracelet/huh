# Huh?

A simple and powerful library for building interactive forms in the terminal. Powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea).

<img alt="Running a burger form" width="600" src="https://vhs.charm.sh/vhs-717b54Ag22l5YUjEgVtpdS.gif">

The above example is running from a single Go program ([source](./examples/burger/main.go)).

## Tutorial

`huh?` provides a straightforward API to build forms and prompt users for input.

For this tutorial, we're building a Burger order form. Lets start by importing
the dependencies that we'll need.

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
    // Prompt the user to choose a burger.
    huh.NewGroup(
        huh.NewSelect[string]().
            Options(
                huh.NewOption("Charmburger Classic", "classic"),
                huh.NewOption("Chickwich", "chickwich"),
                huh.NewOption("Fishburger", "Fishburger"),
                huh.NewOption("Charmpossible™ Burger", "charmpossible"),
            ).
            Title("Choose your burger").
            Value(&burger),
    ),

    // Prompt for toppings and special instructions.
    // The customer can ask for up to 4 toppings.
    huh.NewGroup(
        huh.NewMultiSelect[string]().
            Options(
                huh.NewOption("Lettuce", "Lettuce").Selected(true),
                huh.NewOption("Tomatoes", "Tomatoes").Selected(true),
                huh.NewOption("Charm Sauce", "Charm Sauce"),
                huh.NewOption("Jalapeños", "Jalapeños"),
                huh.NewOption("Cheese", "Cheese"),
                huh.NewOption("Vegan Cheese", "Vegan Cheese"),
                huh.NewOption("Nutella", "Nutella"),
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

<img alt="Input field" width="600" src="https://vhs.charm.sh/vhs-1ULe9JbTHfwFmm3hweRVtD.gif">

```go
huh.NewInput().
    Title("What's for lunch?").
    Prompt("?").
    Validate(isFood).
    Value(&lunch)
```

### Text

<img alt="Text field" width="600" src="https://vhs.charm.sh/vhs-2rrIuVSEf38bT0cwc8hfEG.gif">

```go
huh.NewText().
    Title("Tell me a story.").
    Validate(checkForPlagiarism).
    Value(&story)
```

### Select

<img alt="Select field" width="600" src="https://vhs.charm.sh/vhs-7wFqZlxMWgbWmOIpBqXJTi.gif">

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

<img alt="Multiselect field" width="600" src="https://vhs.charm.sh/vhs-7bYeyKzQNGPyPdpQqvFwkx.gif">

```go
huh.NewMultiSelect[string]().
    Options(
        huh.NewOption("Lettuce", "Lettuce").Selected(true),
        huh.NewOption("Tomatoes", "Tomatoes").Selected(true),
        huh.NewOption("Charm Sauce", "Charm Sauce"),
        huh.NewOption("Jalapeños", "Jalapeños"),
        huh.NewOption("Cheese", "Cheese"),
        huh.NewOption("Vegan Cheese", "Vegan Cheese"),
        huh.NewOption("Nutella", "Nutella"),
    ).
    Title("Toppings").
    Limit(4).
    Value(&toppings)
```

### Confirm

<img alt="Confirm field" width="600" src="https://vhs.charm.sh/vhs-3pCtJgM9EH1tcO0VNtSV3I.gif">

```go
huh.NewConfirm().
    Title("You sure?").
    Affirmative("Yes!").
    Negative("No.").
    Value(&confirm)
```

## Accessibility

Forms can be made accessible to screen readers through setting the
`WithAccessible` option. It's useful to set this through an environment variable
or configuration option to allow the user to control whether their form is
accessible for screen readers.

```go
form.WithAccessible(os.Getenv("ACCESSIBLE") != "")
```

Making the form accessible will remove redrawing and use more standard prompts
to ensure that screen readers are able to dictate the information on screen
correctly.

<img alt="Accessible cuisine form" width="600" src="https://vhs.charm.sh/vhs-19xEBn4LgzPZDtgzXRRJYS.gif">

## Themes

Forms can be customized through themes. You can supply your own custom theme
or use the predefined themes.

There are currently four predefined themes:

* `Charm`
* `Dracula`
* `Base 16`
* `Default`

<br />
<p>
    <img alt="Charm-themed form" width="400" src="https://stuff.charm.sh/huh/themes/charm-theme.png">
    <img alt="Dracula-themed form" width="400" src="https://stuff.charm.sh/huh/themes/dracula-theme.png">
    <img alt="Base 16-themed form" width="400" src="https://stuff.charm.sh/huh/themes/basesixteen-theme.png">
    <img alt="Default-themed form" width="400" src="https://stuff.charm.sh/huh/themes/default-theme.png">
</p>

## Spinner

Huh additionally provides a `spinner` subpackage for displaying spinners while
performing actions. It's useful to complete an action after your user completes
a form.

<img alt="Spinner while making a burger" width="600" src="https://vhs.charm.sh/vhs-5uVCseHk9F5C4MdtZdwhIc.gif">

To get started, create a new spinner, set a title, set an action, and run the
spinner:

```go
makeBurger := func() {
    //...
}

err := spinner.New().
    Title("Making your burger...").
    Action(makeBurger).
    Run()

fmt.Println("Order up!")
```

Alternatively, you can also use `Context`s. The spinner will stop once the
context is cancelled.

```go
makeBurger := func() {
    // ...
}

ctx, _ := context.WithTimeout(context.Background(), time.Second)

go makeBurger()

err := spinner.New().
    Title("Making your burger...").
    Context(ctx).
    Run()

fmt.Println("Order up!")
```

## What about [Bubble Tea][tea]?

Huh doesn’t replace Bubble Tea. Rather, it is an abstraction built on Bubble Tea
to make forms easier to code and implement. It was designed to make assembling
powerful and feature-rich forms in Go as simple and fun as possible.

While you can use `huh` as a replacement to Bubble Tea in many applications
where you only need to prompt the user for input. You can embed `huh` forms in
Bubble Tea applications and use the form as a [Bubble][bubbles].

```go
type Model struct {
    // embed form in parent model, use as bubble.
    form *huh.Form
}
```

See [the Bubble Tea example][example] for how to embed `huh` forms in [Bubble
Tea][tea] applications for more advanced use cases.

[tea]: https://github.com/charmbracelet/bubbletea
[bubbles]: https://github.com/charmbracelet/bubbles
[example]: https://github.com/charmbracelet/huh/blob/main/examples/bubbletea/main.go


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
