package main

import (
	"atlas.sql/internal/db"
	"atlas.sql/internal/ui"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var Version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("atlas.sql v%s\n", Version)
		return
	}

	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		showHelp()
		return
	}

	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	connStr := os.Args[1]
	database, err := db.Connect(connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	p := tea.NewProgram(ui.NewModel(database), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("Atlas SQL - Terminal-based SQL client")
	fmt.Println("\nUsage:")
	fmt.Println("  atlas.sql [connection_string]")
	fmt.Println("\nExamples:")
	fmt.Println("  SQLite:     atlas.sql sqlite://path/to/db.sqlite")
	fmt.Println("  PostgreSQL: atlas.sql \"postgres://user:pass@localhost:5432/dbname?sslmode=disable\"")
	fmt.Println("\nControls:")
	fmt.Println("  Enter: Run Query")
	fmt.Println("  Tab: Switch focus")
	fmt.Println("  Arrows: Navigate Results / Columns")
	fmt.Println("  +/-: Adjust Column Width")
	fmt.Println("  v: Detail View")
	fmt.Println("  c: Copy CSV")
	fmt.Println("  h: Show Help")
	fmt.Println("  q: Quit")
}
