package faq

import "github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"

type FAQService struct {
	repo *FAQRepository
}

func NewFAQService(repo *FAQRepository) *FAQService {
	return &FAQService{repo: repo}
}
func (s *FAQService) GetAll() ([]models.FAQResponce, error) {
	faqs, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	responce := make([]models.FAQResponce, 0, len(faqs))

	for _, f := range faqs {
		responce = append(responce, models.FAQResponce{
			Question: f.Question,
			Answer:   f.Answer,
		})
	}
	return responce, nil
}
