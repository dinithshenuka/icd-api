package icd11

// ICDCode represents the domain model for an ICD classification
type ICDCode struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Version     string `json:"version"`
}
