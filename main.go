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
	var connStr string
	var listTables bool

	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	for _, arg := range os.Args[1:] {
		if arg == "-v" || arg == "--version" {
			fmt.Printf("atlas.sql v%s\n", Version)
			return
		} else if arg == "-h" || arg == "--help" {
			showHelp()
			return
		} else if arg == "-l" || arg == "--list" || arg == "--tables" {
			listTables = true
		} else if connStr == "" {
			connStr = arg
		}
	}

	if connStr == "" {
		showHelp()
		os.Exit(1)
	}

	database, err := db.Connect(connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	if listTables {
		result, err := database.ListTables()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing tables: %v\n", err)
			os.Exit(1)
		}
		if len(result.Rows) == 0 {
			fmt.Println("No tables found.")
		} else {
			for _, row := range result.Rows {
				if len(row) > 0 {
					fmt.Println(row[0])
				}
			}
		}
		return
	}

	p := tea.NewProgram(ui.NewModel(database), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("Atlas SQL - Terminal-based SQL client")
	fmt.Println("\nUsage:")
	fmt.Println("  atlas.sql [options] [connection_string]")
	fmt.Println("\nOptions:")
	fmt.Println("  -l, --list    List all tables in the database and exit")
	fmt.Println("  -v, --version Show version and exit")
	fmt.Println("  -h, --help    Show this help message and exit")
	fmt.Println("\nExamples:")
	fmt.Println("  SQLite:     atlas.sql sqlite://path/to/db.sqlite")
	fmt.Println("  PostgreSQL: atlas.sql \"postgres://user:pass@localhost:5432/dbname?sslmode=disable\"")
	fmt.Println("  List Table: atlas.sql -l sqlite://path/to/db.sqlite")
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
