package res_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/res"
)

func TestJsonResponse(t *testing.T) {

	rec := httptest.NewRecorder()
	data := func() {}
	res.JsonResponse(rec, data, 200)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("ожидали 500, получили %d", rec.Code)
	}

	body := strings.TrimSpace(rec.Body.String())

	expected := `{"error":"ошибка: не удалось преобразовать JSON"}`

	if body != expected {
		t.Errorf("ожидали тело '%s', получили '%s'", expected, body)
	}
}
