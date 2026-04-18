package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iammrdp/icd-api/internal/icdcodes"
)

func main() {
	// Initialize layers
	repo := icdcodes.NewRepository()
	service := icdcodes.NewService(repo)
	handler := icdcodes.NewHandler(service)

	r := gin.Default()

	// API Versioning using Groups
	v1 := r.Group("/v1")
	{
		v1.GET("/icd10", handler.GetICD10)
		v1.GET("/search", handler.Search)
	}

	r.Run(":8080")
}
