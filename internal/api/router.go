package api

import (
	"github.com/gin-gonic/gin"
	"github.com/iammrdp/icd-api/internal/icdcodes"
)

// NewRouter sets up the entire API routing
func NewRouter(handler icdcodes.ServerInterface) *gin.Engine {
	r := gin.Default()

	// 1. Documentation routes
	r.StaticFile("/openapi.yaml", "./api/v1/openapi.yaml")
	r.GET("/docs", ServeSwaggerUI)

	// 2. API v1 logic
	v1 := r.Group("/v1")
	icdcodes.RegisterHandlers(v1, handler)

	return r
}
