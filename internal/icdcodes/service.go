package icdcodes

// Service handles the business logic
type Service struct {
	repo *Repository
}

// NewService creates a new service instance
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetCodes returns the list of codes from the repository
func (s *Service) GetCodes() []ICDCode {
	// You could add logic here: logging, filtering, etc.
	return s.repo.GetAll()
}
