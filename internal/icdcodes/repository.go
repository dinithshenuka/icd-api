package icdcodes

// Repository handles data storage logic
type Repository struct {
	// We could put a database connection here later
}

// NewRepository creates a new instance of the repository
func NewRepository() *Repository {
	return &Repository{}
}

// GetAll returns a hardcoded list of ICD codes
func (r *Repository) GetAll() []ICDCode {
	return []ICDCode{
		{ID: "1", Code: "A00", Description: "Cholera", Version: "ICD-10"},
		{ID: "2", Code: "B00", Description: "Herpesviral infections", Version: "ICD-10"},
	}
}

// Search filters codes by description or code
func (r *Repository) Search(query string) []ICDCode {
	all := r.GetAll()
	var results []ICDCode

	// Basic filter logic (case-insensitive search would be better, but we'll start here)
	for _, item := range all {
		if query == "" || item.Code == query || item.Description == query {
			results = append(results, item)
		}
	}
	return results
}
