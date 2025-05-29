package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/huh"
	// "github.com/charmbracelet/lipgloss" // Not used in the provided example, can be removed
)

func main() {
	var selectedRow any // Or more specific type if known, e.g., table.Row

	// Define columns for the table
	cols := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Name", Width: 15},
		{Title: "Age", Width: 5},
		{Title: "City", Width: 20},
	}

	// Define sample data for the table rows
	rows := []table.Row{
		{"1", "Alice", "30", "New York"},
		{"2", "Bob", "24", "Los Angeles"},
		{"3", "Charlie", "36", "Chicago"},
		{"4", "Diana", "28", "Houston"},
		{"5", "Edward", "45", "Phoenix"},
	}

	// Create the table field
	tableField := huh.NewTable().
		Key("userInfo"). // Key for retrieving the value
		Title("Select User Information").
		Description("Please select a user from the table below.").
		Columns(cols).
		Rows(rows).
		Height(10).          // Set a height for the table viewport
		Value(&selectedRow) // Bind the selected row to the variable

	// Add the table field to a form
	form := huh.NewForm(
		huh.NewGroup(tableField),
	)

	// Run the form
	fmt.Println("--- Running Table in Interactive Mode ---")
	err := form.Run()
	if err != nil {
		// huh.ErrUserAborted is a common error to check for.
		if err == huh.ErrUserAborted {
			fmt.Println("Form aborted by user.")
		} else {
			log.Fatalf("Form error: %v", err)
		}
	}

	// Print the selected data
	if selectedRow != nil {
		// Assuming selectedRow is of type table.Row or similar that can be printed
		// If it's `any`, you might need a type assertion
		if sr, ok := selectedRow.(table.Row); ok {
			fmt.Printf("Interactive Mode - Selected: ID=%s, Name=%s, Age=%s, City=%s\n", sr[0], sr[1], sr[2], sr[3])
		} else {
			fmt.Printf("Interactive Mode - Selected row data: %v\n", selectedRow)
		}
	} else {
		fmt.Println("Interactive Mode: No row selected or form was aborted.")
	}

	// Example of how to use the table in accessible mode
	// Create another table field for demonstration
	var accessibleSelectedRow any
	accessibleTableField := huh.NewTable().
		Key("accessibleUser").
		Title("Select User (Accessible Mode)").
		Columns(cols).
		Rows(rows).
		Height(5). // Will be ignored in accessible mode but good practice
		Value(&accessibleSelectedRow).
		WithAccessible(true) // Enable accessible mode

	fmt.Println("\n--- Running Table in Accessible Mode ---")
	// To run a single field in accessible mode directly (without a form):
	// The Run() method on the field itself handles accessible mode.
	err = accessibleTableField.Run()
	if err != nil {
		if err == huh.ErrUserAborted {
			fmt.Println("Accessible table input aborted by user.")
		} else {
			log.Fatalf("Accessible table error: %v", err)
		}
	}

	if accessibleSelectedRow != nil {
		if sr, ok := accessibleSelectedRow.(table.Row); ok {
			fmt.Printf("Accessible Mode - Selected: ID=%s, Name=%s, Age=%s, City=%s\n", sr[0], sr[1], sr[2], sr[3])
		} else {
			fmt.Printf("Accessible Mode - Selected row data: %v\n", accessibleSelectedRow)
		}
	} else {
		fmt.Println("Accessible table: No row selected or input aborted.")
	}
}
