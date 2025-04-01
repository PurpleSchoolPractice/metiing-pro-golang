package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/jwt"
)

type key string

const (
	ContextEmailKey key = "ContextEmailKey"
)

func writeUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	_, err := w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
	if err != nil {
		return
	}
}

func IsAuthed(next http.Handler, jwtService *jwt.JWT) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" || !strings.HasPrefix(authorization, "Bearer ") {
			writeUnauthorized(w)
			return
		}

		token := strings.TrimPrefix(authorization, "Bearer ")
		isValid, data := jwtService.ParseToken(token)
		if !isValid || data == nil {
			writeUnauthorized(w)
			return
		}

		ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
