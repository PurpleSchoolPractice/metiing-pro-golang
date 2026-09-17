package models

import "gorm.io/gorm"

type FAQ struct {
	gorm.Model
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
type FAQResponce struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
