package main

import (
	"log"
	"os"
	"time"

	"github.com/iammrdp/icd-api/internal/api"
	"github.com/iammrdp/icd-api/internal/api/handler"
	"github.com/iammrdp/icd-api/internal/domain/icd11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {
	// 1. Setup Structured Logs
	zerolog.TimeFieldFormat = time.RFC3339
	if os.Getenv("GIN_MODE") != "release" {
		log.SetOutput(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen})
	}

	// 2. Load environment
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "database/icd11.db"
	}

	// 3. Initialize Domain Layers
	repo := icd11.NewRepository(dbPath)
	service := icd11.NewService(repo)
	
	// 4. Initialize Delivery Layer
	h := handler.NewICD11Handler(service)

	// 5. Setup and Start Server
	r := api.NewRouter(h)
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
