package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/dinithshenuka/icd-code-api/internal/domain/icd11"
)

type ICD11Handler struct {
	service *icd11.Service
}

func NewICD11Handler(s *icd11.Service) *ICD11Handler {
	return &ICD11Handler{service: s}
}

func (h *ICD11Handler) GetIcd10(c *gin.Context) {
	codes := h.service.GetCodes()
	c.JSON(http.StatusOK, codes)
}

func (h *ICD11Handler) SearchCodes(c *gin.Context, params SearchCodesParams) {
	// Safely check for nil pointer or empty string
	if params.Q == nil || *params.Q == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Query parameter 'q' is required",
		})
		return
	}

	results := h.service.SearchCodes(*params.Q)
	c.JSON(http.StatusOK, results)
}
