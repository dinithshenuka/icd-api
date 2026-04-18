package main

import (
	"github.com/iammrdp/icd-api/internal/api"
	"github.com/iammrdp/icd-api/internal/icdcodes"
)

func main() {
	// 1. Initialize logic layers
	repo := icdcodes.NewRepository()
	service := icdcodes.NewService(repo)
	handler := icdcodes.NewHandler(service)

	// 2. Setup the router from our internal api package
	r := api.NewRouter(handler)

	// 3. Start the engine
	r.Run(":8080")
}
