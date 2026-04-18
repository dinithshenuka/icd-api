package icd11

// Service handles business logic for ICD codes
type Service struct {
	repo *Repository
}

// NewService creates a new ICD code service
func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) GetCodes() []ICDCode {
	return s.repo.GetAll()
}

func (s *Service) SearchCodes(query string) []ICDCode {
	return s.repo.Search(query)
}
