package main

import (
	"github.com/iammrdp/icd-api/internal/api"
	"github.com/iammrdp/icd-api/internal/icdcodes"
)

func main() {
	// 1. Initialize logic layers with the real database
	repo := icdcodes.NewRepository("database/icd11.db")
	service := icdcodes.NewService(repo)
	handler := icdcodes.NewHandler(service)

	// 2. Setup the router
	r := api.NewRouter(handler)

	// 3. Start the engine
	r.Run(":8080")
}
