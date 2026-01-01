package main

import (
	"context" // Import context
	"log"
	"net/http" // Import net/http
	"os"       // Import os
	"os/signal" // Import os/signal
	"syscall"  // Import syscall
	"time"     // Import time

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/smartq/smartq/internal/api"
	"github.com/smartq/smartq/internal/config"
	"github.com/smartq/smartq/internal/notifier" // Import the notifier package
	"github.com/smartq/smartq/internal/storage"
)

func main() {
	cfg, err := config.LoadConfig() // Load config first
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Run database migrations
	runMigrations(cfg.DatabaseURL)

	db, err := storage.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create and run the WebSocket hub
	hub := notifier.NewHub()
	go hub.Run()

	router := api.NewRouter(db, hub) // Pass the hub to the router

	// Start HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port, // Use config.Port
		Handler: router,
	}

	// Start HTTP server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
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

