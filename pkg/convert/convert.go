package convert

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// парсим ID в строке
func ParseId(r *http.Request, parsIs string) (uint, error) {
	strUserID := chi.URLParam(r, parsIs)

	idUint, err := strconv.ParseUint(strUserID, 10, 64)
	if err != nil {

		return 0, errors.New("not a valid ID")
	}
	id := uint(idUint)
	return id, nil
}
