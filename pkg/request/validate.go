package request

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

func Validate[T any](payload T) error {
	validate := validator.New()
	err := validate.Struct(&payload)
	return err
}

// парсим время из запроса и валидируем
func ValidateTime(date string) (time.Time, error) {
	layout := "2006-01-02 15:04"
	var zeroTime time.Time
	parseTime, err := time.Parse(layout, date)
	if err != nil {
		return zeroTime, errors.New("wrong time in request. Format time shoud be 2006-01-02 15:04")
	}

	return parseTime, nil
}
