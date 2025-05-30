package convert

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// преобразуем ID в UINT
<<<<<<< HEAD
func ParseId(r *http.Request) (uint, error) {
=======
func ConvertID(r *http.Request) (uint, error) {
>>>>>>> master
	strUserID := chi.URLParam(r, "id")

	idUint, err := strconv.ParseUint(strUserID, 10, 64)
	if err != nil {

		return 0, errors.New("not a valid ID")
	}
	id := uint(idUint)
	return id, nil
}
