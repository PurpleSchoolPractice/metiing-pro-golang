package convert

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

var ErrInvalidID = errors.New("not a valid ID")

// парсим ID в строке
func ParseId(r *http.Request, parsIs string) (uint, error) {
	strUserID := chi.URLParam(r, parsIs)

	parsedID, err := strconv.ParseUint(strUserID, 10, strconv.IntSize)
	if err != nil || parsedID == 0 {
    	return 0, ErrInvalidID
    }
	id := uint(parsedID)
	return id, nil
}
