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

		value := r.Context().Value(middleware.ContextEmailKey)
		emailKey := value.(string)
		value = r.Context().Value(middleware.ContextUserIDKey)
		userIDKey := value.(uint)

		if emailKey != "test@example.com" || userIDKey != 42 {
			t.Error("данные в контексте и в токене не совпадают")
		}
	})

	middleware.IsAuthed(next, newJWT).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("ожидали 200, получили %d", rr.Code)
	}

}

func TestIsAuthedNegative(t *testing.T) {
	testCases := []struct {
		name   string
		header string
		status int
	}{
		{name: "No header", header: "", status: 401},
		{name: "Wrong prefix", header: "Token xxx", status: 401},
		{name: "Crash token", header: "Bearer invalid", status: 401},
	}

	newJWT := jwt.NewJWT("secret")
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("дошли до хендлера — а не должны были!")
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com/", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			rr := httptest.NewRecorder()
			middleware.IsAuthed(next, newJWT).ServeHTTP(rr, req)
			if rr.Code != 401 {
				t.Errorf("ожидали 401, получили %d", rr.Code)
			}
		})
	}
}
