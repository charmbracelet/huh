package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type Category string

const (
	CategoryFood     Category = "food"
	CategoryClothes  Category = "clothes"
	CategoryTools    Category = "tool"
	CategoryCleaning Category = "cleaning"
)

var (
	CategoryLabels = map[Category]string{
		CategoryFood:     "Food",
		CategoryCleaning: "Cleaning",
		CategoryClothes:  "Clothes",
		CategoryTools:    "Tools",
	}
)

type Item struct {
	Name     string
	Category Category
	Price    uint
}

var (
	Items = []Item{
		{
			Name:     "Banana",
			Category: CategoryFood,
			Price:    1,
		},
		{
			Name:     "Nem chua",
			Category: CategoryFood,
			Price:    3,
		},
		{
			Name:     "Lemonade",
			Category: CategoryFood,
			Price:    5,
		},
		{
			Name:     "Shirt",
			Category: CategoryClothes,
			Price:    15,
		},
		{
			Name:     "Trousers",
			Category: CategoryClothes,
			Price:    60,
		},
		{
			Name:     "Vest",
			Category: CategoryClothes,
			Price:    40,
		},
		{
			Name:     "Sponge",
			Category: CategoryCleaning,
			Price:    8,
		},
		{
			Name:     "Detergent",
			Category: CategoryCleaning,
			Price:    4,
		},
		{
			Name:     "Broom",
			Category: CategoryCleaning,
			Price:    10,
		},
		{
			Name:     "Nails",
			Category: CategoryTools,
			Price:    12,
		},
		{
			Name:     "Hammer",
			Category: CategoryTools,
			Price:    40,
		},
		{
			Name:     "Saw",
			Category: CategoryTools,
			Price:    30,
		},
	}
)

func main() {
	columns := []huh.Column[Item]{
		huh.NewColumn("Item", 25, func(i Item) any { return i.Name }),
		huh.NewColumn("Category", 25, func(i Item) any { return CategoryLabels[i.Category] }),
		huh.NewColumn("Price", 10, func(i Item) any { return fmt.Sprintf("%9d€", i.Price) }),
	}

	options := huh.NewTableOptions(func(i Item) string { return i.Name }, Items...)

	theme := huh.ThemeCharm()
	theme.Focused.Table.Header = theme.Focused.Table.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	theme.Blurred.Table.Header = theme.Blurred.Table.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	theme.Focused.Table.Selected = theme.Focused.Table.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	theme.Blurred.Table.Selected = theme.Focused.Table.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	table := huh.NewTable[Item, string]().
		Title("Pick a product").
		Description("Enjoy our large selection").
		Height(8).
		Width(60).
		Columns(columns...).
		Options(options...).
		WithTheme(theme)

	if err := table.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.GetValue())
}
