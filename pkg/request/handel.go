package request

import (
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/res"
	"net/http"
)

func HandelBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		res.JsonResponse(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	err = Validate[T](body)
	if err != nil {
		res.JsonResponse(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return &body, nil
}
