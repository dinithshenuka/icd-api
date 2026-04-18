package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iammrdp/icd-api/internal/icdcodes"
)

func main() {
	// 1. Initialize layers
	repo := icdcodes.NewRepository()
	service := icdcodes.NewService(repo)
	handler := icdcodes.NewHandler(service)

	// 2. Set up router
	r := gin.Default()

	// 3. Define routes
	r.GET("/icd10", handler.GetICD10)

	// 4. Start server
	r.Run(":8080")
}
