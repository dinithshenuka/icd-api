package icdcodes

// ICDCode represents a clinical classification code
type ICDCode struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Version     string `json:"version"` // e.g., ICD-10 or ICD-11
}
