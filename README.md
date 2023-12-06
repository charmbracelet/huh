# Huh?

A simple, powerful library for building interactive forms in the terminal.
Powered by [Bubble Tea][tea].

<img alt="Running a burger form" width="600" src="https://vhs.charm.sh/vhs-3J4i6HE3yBmz6SUO3HqILr.gif">

The above example is running from a single Go program ([source](./examples/burger/main.go)).

## Tutorial

`huh?` is a Go library for prompting users for input. You can build complex
forms in a few lines of Go.

Let’s build a Burger order form. To start, let's import `charmbracelet/huh` and
define a few variables to store the data we'll prompt for.

```go
package main

import (
  "log"

  "github.com/charmbracelet/huh"
)

var (
    burger string
    toppings []string
    name string
    instructions string
    discount bool
)
```

`huh` separates forms into groups (you can think of groups as pages). Groups
are made of fields (e.g. `Select`, `Input`, `Text`). We will set up three
groups for the customer to fill out.

```go
form := huh.NewForm(
    // Ask the user for a base burger and toppings.
    huh.NewGroup(
        huh.NewSelect[string]().
            Options(
                huh.NewOption("Charmburger Classic", "classic"),
                huh.NewOption("Chickwich", "chickwich"),
                huh.NewOption("Fishburger", "fishburger"),
                huh.NewOption("Charmpossible™ Burger", "charmpossible"),
            ).
            Title("Choose your burger").
            Value(&burger),
    ),

    // Let the user select multiple toppings.
    // We allow a maximum limit of 4 toppings.
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

    // Gather some final details about the order.
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
            Title("Would you like 15% off?").
            Value(&discount),
    ),
)
```

Finally, run the form:

```go
err := form.Run()
if err != nil {
    log.Fatal(err)
}

if !discount {
    fmt.Println("What?? You didn’t take the discount?!"
}
```

And that’s it! For more info see [the full source
code](./examples/burger/main.go) for this example and
[the docs](https://pkg.go.dev/github.com/charmbracelet/huh?tab=doc).

## Field Reference

- [`Input`](#input): single line text input
- [`Text`](#text): multi-line text input
- [`Select`](#select): select an option from a list
- [`MultiSelect`](#multiple-select): select multiple options from a list
- [`Confirm`](#confirm): confirm an action (yes or no)

> [!TIP]
> Just want to prompt the user with a single field? Each field has a `Run`
> method that can be used as a shorthand for gathering quick and easy input.

```go
var name string

huh.NewInput().
    Title("What's your name?").
    Value(&name).
    Run() // this is blocking...

fmt.Printf("Hey, %s!\n", name)
```

### Input

Prompt the user for a single line of text.

<img alt="Input field" width="600" src="https://vhs.charm.sh/vhs-1ULe9JbTHfwFmm3hweRVtD.gif">

```go
huh.NewInput().
    Title("What's for lunch?").
    Prompt("?").
    Validate(isFood).
    Value(&lunch)
```

### Text

Prompt the user for multiple lines of text.

<img alt="Text field" width="600" src="https://vhs.charm.sh/vhs-2rrIuVSEf38bT0cwc8hfEG.gif">

```go
huh.NewText().
    Title("Tell me a story.").
    Validate(checkForPlagiarism).
    Value(&story)
```

### Select

Prompt the user to select a single option from a list.

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

Prompt the user to select multiple (zero or more) options from a list.

<img alt="Multiselect field" width="600" src="https://vhs.charm.sh/vhs-3TLImcoexOehRNLELysMpK.gif">

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

Prompt the user to confirm (Yes or No).

<img alt="Confirm field" width="600" src="https://vhs.charm.sh/vhs-2HeX5MdOxLsrWwsa0TNMIL.gif">

```go
huh.NewConfirm().
    Title("You sure?").
    Affirmative("Yes!").
    Negative("No.").
    Value(&confirm)
```

## Accessibility

Prevent redrawing the screen with the `WithAccessible` option. This is useful
for screen readers to provide better dictation of the output.

> [!TIP]
> We recommend setting this through an environment variable or configuration
> option to allow the user to control accessibility.

```go
accessibleMode := os.Getenv("ACCESSIBLE") != ""
form.WithAccessible(accessibleMode)
```

Accessible forms will remove redrawing in favor of standard prompts, providing
better dictation of the information on screen for the visually impaired.

<img alt="Accessible cuisine form" width="600" src="https://vhs.charm.sh/vhs-19xEBn4LgzPZDtgzXRRJYS.gif">

## Themes

Forms can be themed.

Supply your own custom theme or choose from one of the four predefined themes:

- `Charm`
- `Dracula`
- `Base 16`
- `Default`

<br />
<p>
    <img alt="Charm-themed form" width="400" src="https://stuff.charm.sh/huh/themes/charm-theme.png">
    <img alt="Dracula-themed form" width="400" src="https://stuff.charm.sh/huh/themes/dracula-theme.png">
    <img alt="Base 16-themed form" width="400" src="https://stuff.charm.sh/huh/themes/basesixteen-theme.png">
    <img alt="Default-themed form" width="400" src="https://stuff.charm.sh/huh/themes/default-theme.png">
</p>

## Spinner

Spinners come built in to `huh` for loading actions. It's useful to indicate
loading while completing an action after a form is submitted.

<img alt="Spinner while making a burger" width="600" src="https://vhs.charm.sh/vhs-5uVCseHk9F5C4MdtZdwhIc.gif">

Create a new spinner, set a title, set the action (or provide a `Context`), and run the spinner:

<table>

<tr>
<td> <strong>Action Style</strong> </td><td> <strong>Context Style</strong> </td></tr>
<tr>
<td>

```go
err := spinner.New().
    Title("Making your burger...").
    Action(makeBurger).
    Run()

fmt.Println("Order up!")
```

</td>
<td>

```go
go makeBurger()

err := spinner.New().
    Title("Making your burger...").
    Context(ctx).
    Run()

fmt.Println("Order up!")
```

</td>
</tr>
</table>

## What about [Bubble Tea][tea]?

Huh is an abstraction built on Bubble Tea to make forms easier to code and
implement. You can use `huh` to replace Bubble Tea if you only need to gather
input from the user via a form.

For more complex use cases, however, you can embed `huh` forms within Bubble Tea
applications to add forms to your applications.

<img alt="Bubbletea + Huh?" width="174" src="https://stuff.charm.sh/huh/bubbletea-huh.png">

```go
type Model struct {
    // embed form in parent model, use as bubble.
    form *huh.Form
}

func NewModel() Model {
    return Model{
        form: huh.NewForm(
            huh.NewGroup(
                huh.NewSelect[string]().
                    Key("class").
                    Options(huh.NewOptions("Warrior", "Mage", "Rogue")...).
                    Title("Choose your class"),

            huh.NewSelect[int]().
                Key("level").
                Options(huh.NewOptions(1, 20, 9999)...).
                Title("Choose your level"),
            ),
        )
    }
}

func (m Model) Init() tea.Cmd {
    return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // ...

    form, cmd := m.form.Update(msg)
    if f, ok := form.(*huh.Form); ok {
        m.form = f
    }

    return m, cmd
}

func (m Model) View() string {
    if m.form.State == huh.StateCompleted {
        class := m.form.GetString("class")
        level := m.form.GetString("level")
        return fmt.Sprintf("You selected: %s, Lvl. %d", class, level)
    }
    return m.form.View()
}

```

<img alt="Bubble Tea embedded form example" width="800" src="https://vhs.charm.sh/vhs-3wGaB7EUKWmojeaHpARMUv.gif">

See [the Bubble Tea example][example] for how to embed `huh` forms in [Bubble
Tea][tea] applications for more advanced use cases.

[tea]: https://github.com/charmbracelet/bubbletea
[bubbles]: https://github.com/charmbracelet/bubbles
[example]: https://github.com/charmbracelet/huh/blob/main/examples/bubbletea/main.go

## Feedback

We'd love to hear your thoughts on this project. Feel free to drop us a note!

- [Twitter](https://twitter.com/charmcli)
- [The Fediverse](https://mastodon.social/@charmcli)
- [Discord](https://charm.sh/chat)

## Acknowledgments

`huh` is inspired by the wonderful [Survey][survey] library by Alec Aivazis.

[survey]: https://github.com/AlecAivazis/survey

## License

[MIT](https://github.com/charmbracelet/bubbletea/raw/master/LICENSE)

---

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400"></a>

Charm热爱开源 • Charm loves open source • نحنُ نحب المصادر المفتوحة
