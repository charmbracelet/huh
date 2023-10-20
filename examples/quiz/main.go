package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
)

type Language int

const (
	None Language = iota
	C
	Erlang
	Go
	Haskell
	Python
	Rust
	Swift
	Typescript
	Zig
)

func (l Language) String() string {
	return map[Language]string{
		None:       "None",
		C:          "C Programmer",
		Erlang:     "Erlang person",
		Go:         "Gopher",
		Haskell:    "Haskeller",
		Python:     "Pythonista",
		Rust:       "Rustacean",
		Swift:      "Swift Lover",
		Typescript: "Typescripter",
		Zig:        "Zig Programmer",
	}[l]
}

type Answer struct {
	Language Language
	Points   int
}

type Quiz struct {
	office Answer
	malloc Answer
	notes  string
}

func main() {
	fmt.Printf("What kind of coder are you?\n\n")

	var q Quiz

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[Answer]().Value(&q.office).
				Title("Where’s your dream office?").
				Required(true).
				Options(
					huh.NewOption("In bed", Answer{Zig, 1}),
					huh.NewOption("Any ol’ desk with a Moonlander", Answer{Rust, 1}),
					huh.NewOption("In Rob Pike’s state machine", Answer{Go, 1}),
					huh.NewOption("Microsoft Office", Answer{Typescript, 1}),
					huh.NewOption("Sony Ericsson corporate headquarters", Answer{Erlang, 1}),
					huh.NewOption("An ivory tower", Answer{Haskell, 1}),
				),
			huh.NewSelect[Answer]().Value(&q.malloc).
				Title("If you don’t malloc() are you even living?").
				Required(true).
				Options(
					huh.NewOption("No, you’re not!", Answer{C, 2}),
					huh.NewOption("Ugh, I have work to do", Answer{Go, 1}),
					huh.NewOption("What?", Answer{Typescript, 1}),
				),
			huh.NewSelect[Answer]().Value(&q.malloc).
				Title("Side effects are:").
				Required(true).
				Options(
					huh.NewOption("Impure", Answer{Rust, 1}),
					huh.NewOption("The best", Answer{Haskell, 1}),
					huh.NewOption("Great when partying", Answer{Haskell, 1}),
				),
		),
		huh.NewGroup(
			huh.NewSelect[Answer]().Value(&q.malloc).
				Title("How do you live your life?").
				Options(
					huh.NewOption("Dangerously", Answer{C, 1}),
					huh.NewOption("Snake-style", Answer{Python, 1}),
					huh.NewOption("IDK", Answer{Go, 1}),
				),
		),
		huh.NewGroup(
			huh.NewText().Value(&q.notes).Required(true).Title("Is there anything else we should know?"),
		),
	)

	if a, err := strconv.ParseBool(os.Getenv("ACCESSIBLE")); err == nil {
		form.Accessible(a)
	}

	if err := form.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	points := make(map[Language]int)

	for _, v := range []Answer{
		q.office,
		q.malloc,
	} {
		points[v.Language] += v.Points
	}

	var winner Language
	for k, v := range points {
		if v > points[winner] {
			winner = k
		}
	}

	if winner != None {
		fmt.Printf("Condragulations, you are a %s!\n", winner)
	}
}
