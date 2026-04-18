package main

import (
	"log"
	"os"

	"github.com/iammrdp/icd-api/internal/api"
	"github.com/iammrdp/icd-api/internal/icdcodes"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 2. Get configuration from Environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "database/icd11.db"
	}

	// 3. Initialize logic layers
	repo := icdcodes.NewRepository(dbPath)
	service := icdcodes.NewService(repo)
	handler := icdcodes.NewHandler(service)

	// 4. Setup the router
	r := api.NewRouter(handler)

	// 5. Start the engine
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
