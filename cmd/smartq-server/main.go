package main

import (
	"log"

	"github.com/smartq/smartq/internal/api"
	"github.com/smartq/smartq/internal/config"
	"github.com/smartq/smartq/internal/storage"
)

func main() {
	cfg := config.LoadConfig()

	db, err := storage.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	router := api.NewRouter(db)
	router.Run(":8080")
}

