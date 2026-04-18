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
	return s.repo.GetAll()
}

// SearchCodes calls the repository search logic
func (s *Service) SearchCodes(query string) []ICDCode {
	return s.repo.Search(query)
}
