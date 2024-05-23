package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"

	"github.com/charmbracelet/huh"
)

func main() {
	log.SetReportTimestamp(false)

	var (
		country string
		state   string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(huh.NewOptions("United States", "Canada", "Mexico")...).
				Value(&country).
				Title("Country").
				Height(5),
			huh.NewSelect[string]().
				Value(&state).
				Height(8).
				TitleFunc(func() string {
					switch country {
					case "United States":
						return "State"
					case "Canada":
						return "Province"
					default:
						return "Territory"
					}
				}, &country).
				OptionsFunc(func() []huh.Option[string] {
					s := states[country]
					// simulate API call
					time.Sleep(1000 * time.Millisecond)
					return huh.NewOptions(s...)
				}, &country /* only this function when `country` changes */),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s, %s\n", state, country)
}

var states = map[string][]string{
	"Canada": {
		"Alberta",
		"British Columbia",
		"Manitoba",
		"New Brunswick",
		"Newfoundland and Labrador",
		"North West Territories",
		"Nova Scotia",
		"Nunavut",
		"Ontario",
		"Prince Edward Island",
		"Quebec",
		"Saskatchewan",
		"Yukon",
	},
	"Mexico": {
		"Aguascalientes",
		"Baja California",
		"Baja California Sur",
		"Campeche",
		"Chiapas",
		"Chihuahua",
		"Coahuila",
		"Colima",
		"Durango",
		"Guanajuato",
		"Guerrero",
		"Hidalgo",
		"Jalisco",
		"México",
		"Mexico City",
		"Michoacán",
		"Morelos",
		"Nayarit",
		"Nuevo León",
		"Oaxaca",
		"Puebla",
		"Querétaro",
		"Quintana Roo",
		"San Luis Potosí",
		"Sinaloa",
		"Sonora",
		"Tabasco",
		"Tamaulipas",
		"Tlaxcala",
		"Veracruz",
		"Ignacio de la Llave",
		"Yucatán",
		"Zacatecas",
	},
	"United States": {
		"Alabama",
		"Alaska",
		"Arizona",
		"Arkansas",
		"California",
		"Colorado",
		"Connecticut",
		"Delaware",
		"Florida",
		"Georgia",
		"Hawaii",
		"Idaho",
		"Illinois",
		"Indiana",
		"Iowa",
		"Kansas",
		"Kentucky",
		"Louisiana",
		"Maine",
		"Maryland",
		"Massachusetts",
		"Michigan",
		"Minnesota",
		"Mississippi",
		"Missouri",
		"Montana",
		"Nebraska",
		"Nevada",
		"New Hampshire",
		"New Jersey",
		"New Mexico",
		"New York",
		"North Carolina",
		"North Dakota",
		"Ohio",
		"Oklahoma",
		"Oregon",
		"Pennsylvania",
		"Rhode Island",
		"South Carolina",
		"South Dakota",
		"Tennessee",
		"Texas",
		"Utah",
		"Vermont",
		"Virginia",
		"Washington",
		"West Virginia",
		"Wisconsin",
		"Wyoming",
	},
}
