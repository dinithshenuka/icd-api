package handler

// ErrorResponse is the standard error format for the API
type ErrorResponse struct {
	Code    string `json:"code" example:"INTERNAL_ERROR"`
	Message string `json:"message" example:"An unexpected error occurred"`
}
