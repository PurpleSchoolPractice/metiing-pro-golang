package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/jwt"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/middleware"
)

func TestIsAuthed(t *testing.T) {
	newJWT := jwt.NewJWT("secret")

	pair, err := newJWT.GenerateTokenPair(jwt.JWTData{
		UserID: 42,
		Email:  "test@example.com",
	})

	if err != nil {
		t.Fatal(err)
	}

	accessToken := pair.AccessToken

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware.IsAuthed(next, newJWT).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("ожидали 200, получили %d", rr.Code)
	}
}

func TestIsAuthedNegative(t *testing.T) {

}
