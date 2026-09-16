package faq

import "github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"

type FAQService struct {
	repo *FAQRepository
}

func NewFAQService(repo *FAQRepository) *FAQService {
	return &FAQService{repo: repo}
}
func (s *FAQService) GetAll() ([]models.FAQ, error) {
	faqs, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return faqs, nil
}
