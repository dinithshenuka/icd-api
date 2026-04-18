package icdcodes

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests and implements ServerInterface
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetIcd10 implements the generated GetIcd10 interface
func (h *Handler) GetIcd10(c *gin.Context) {
	codes := h.service.GetCodes()
	
	var response []ICDCodeResponse
	for _, code := range codes {
		response = append(response, ICDCodeResponse{
			Code:        code.Code,
			Description: code.Description,
		})
	}
	
	c.JSON(http.StatusOK, response)
}

// SearchCodes implements the generated SearchCodes interface
func (h *Handler) SearchCodes(c *gin.Context, params SearchCodesParams) {
	query := ""
	if params.Q != nil {
		query = *params.Q
	}

	results := h.service.SearchCodes(query)
	
	var response []ICDCodeResponse
	for _, res := range results {
		response = append(response, ICDCodeResponse{
			Code:        res.Code,
			Description: res.Description,
		})
	}
	
	c.JSON(http.StatusOK, response)
}
