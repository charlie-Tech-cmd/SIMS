package main

import (
	"bufio"
	"log"
	"os"
	"sims/storage"
)

func main() {

	// Initialize database connection
	dbEngine, err := storage.NewSQLiteEngine("school.db")
	if err != nil {
		log.Fatalf("could not start application: %v", err)
	}
	defer dbEngine.Close()

	// Create application instance
	sys := &SchoolSystem{
		DB:             dbEngine,
		ConsoleScanner: bufio.NewScanner(os.Stdin),
	}

	// Display startup banner
	sys.DisplaySystemHeader("University Portal")

	// Add default records if the database is empty
	sys.EnsureSeedDataExist()

	// Start the main application menu
	sys.RunMainMenuLoop()

}
