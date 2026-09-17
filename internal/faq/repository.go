package faq

import (
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
)

type FAQRepository struct {
	database *db.Db
}

func NewFAQRepository(database *db.Db) *FAQRepository {
	return &FAQRepository{database: database}
}
func (r *FAQRepository) GetAll() ([]models.FAQ, error) {
	var faqs []models.FAQ
	err := r.database.DB.Find(&faqs).Error
	if err != nil {
		return nil, err
	}
	return faqs, nil
}
