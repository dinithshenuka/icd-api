package icdcodes

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new handler instance
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetICD10 returns all ICD-10 codes formatted as DTOs
func (h *Handler) GetICD10(c *gin.Context) {
	codes := h.service.GetCodes()
	c.JSON(http.StatusOK, MapToResponses(codes))
}

// Search handles the search request and returns DTO results
func (h *Handler) Search(c *gin.Context) {
	query := c.Query("q")
	results := h.service.SearchCodes(query)
	c.JSON(http.StatusOK, MapToResponses(results))
}
