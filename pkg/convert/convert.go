package convert

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// преобразуем ID в UINT
func ConvertID(r *http.Request) (uint, error) {
	strUserID := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		return 0, errors.New("failed to convert")
	}
	return uint(userID), nil
}
