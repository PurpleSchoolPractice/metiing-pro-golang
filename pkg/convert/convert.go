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

	idUint, err := strconv.ParseUint(strUserID, 10, 64)
	if err != nil {

		return 0, errors.New("not a valid ID")
	}
	id := uint(idUint)
	return id, nil
}
