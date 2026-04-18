package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/iammrdp/icd-api/internal/api/handler"
	"github.com/iammrdp/icd-api/internal/api/middleware"
)

// NewRouter sets up the entire API routing
func NewRouter(h *handler.ICD11Handler) *gin.Engine {
	r := gin.New()

	// 1. Global Production Middleware
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	// 2. Documentation & System routes
	r.StaticFile("/openapi.yaml", "./openapi/openapi.yaml")
	r.GET("/docs", ServeSwaggerUI)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "up"})
	})

	// 3. API v1 logic
	v1 := r.Group("/v1")
	handler.RegisterHandlers(v1, h)

	return r
}
