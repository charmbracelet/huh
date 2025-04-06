package main

import "github.com/charmbracelet/huh/v2"

type Pokemon struct {
	id   int
	name string
}

var pokemons = []Pokemon{
	{1, "Bulbasaur"},
	{2, "Ivysaur"},
	{3, "Venusaur"},
	{4, "Charmander"},
	{5, "Charmeleon"},
	{6, "Charizard"},
	{7, "Squirtle"},
	{8, "Wartortle"},
	{9, "Blastoise"},
	{10, "Caterpie"},
	{11, "Metapod"},
	{12, "Butterfree"},
	{13, "Weedle"},
	{14, "Kakuna"},
	{15, "Beedrill"},
	{16, "Pidgey"},
	{17, "Pidgeotto"},
	{18, "Pidgeot"},
	{19, "Rattata"},
	{20, "Raticate"},
	{21, "Spearow"},
	{22, "Fearow"},
	{23, "Ekans"},
	{24, "Arbok"},
	{25, "Pikachu"},
	{26, "Raichu"},
	{27, "Sandshrew"},
	{28, "Sandslash"},
}

func (p Pokemon) String() string {
	return p.name
}

func main() {
	var pokemon Pokemon

	s := huh.NewSelect[Pokemon]().
		Title("Choose your starter").
		Options(huh.NewOptions(pokemons...)...).
		Value(&pokemon).
		WithHeight(7)

	huh.NewForm(huh.NewGroup(s)).Run()
}
