package icdcodes

// ICDCodeResponse is the DTO sent to the API user
type ICDCodeResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// ToResponse converts a internal model to a DTO
func (i *ICDCode) ToResponse() ICDCodeResponse {
	return ICDCodeResponse{
		Code:        i.Code,
		Description: i.Description,
	}
}

// MapToResponses is a helper to convert a slice of models to DTOs
func MapToResponses(codes []ICDCode) []ICDCodeResponse {
	responses := make([]ICDCodeResponse, len(codes))
	for idx, code := range codes {
		responses[idx] = code.ToResponse()
	}
	return responses
}
