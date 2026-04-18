package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var migrationDir = flag.String("path", "database/migrations", "Path to migration files")
	var dbPath = flag.String("db", "database/icd11.db", "Path to SQLite database")
	var action = flag.String("action", "up", "Action to perform (up or down)")

	flag.Parse()

	dbURL := fmt.Sprintf("sqlite://%s", *dbPath)
	m, err := migrate.New(fmt.Sprintf("file://%s", *migrationDir), dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
	}

	switch *action {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration UP failed: %v", err)
		}
		fmt.Println("Successfully applied migrations (UP)")
	case "down":
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration DOWN failed: %v", err)
		}
		fmt.Println("Successfully reverted last migration (DOWN)")
	default:
		log.Fatalf("Unknown action: %s", *action)
	}
}
