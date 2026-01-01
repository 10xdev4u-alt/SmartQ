package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/smartq/smartq/internal/api"
	"github.com/smartq/smartq/internal/config"
	"github.com/smartq/smartq/internal/storage"
)

func main() {
	cfg := config.LoadConfig()

	// Run database migrations
	runMigrations(cfg.DatabaseURL)

	db, err := storage.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	router := api.NewRouter(db)
	router.Run(":8080")
}

func runMigrations(databaseURL string) {
	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	log.Println("Database migrations applied successfully!")
}

