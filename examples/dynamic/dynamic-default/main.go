package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
)

func main() {

	var name, id, directory string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Project name").Value(&name).Placeholder("<Placeholder>").Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Project ID").Value(&id).DefaultFunc(func() string {
				return strings.ReplaceAll(strings.ToLower(name), "/", "-")
			}, &name),
			huh.NewInput().Title("Project directory").Value(&directory).DefaultFunc(func() string {
				if id == "" {
					return "~/projects/<Project ID>"
				}
				return "~/projects/" + id
			}, &id),
		),
	)

	if err := form.Run(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Project Name: ", name)
	fmt.Println("Project ID: ", id)
	fmt.Println("Directory: ", directory)
}
