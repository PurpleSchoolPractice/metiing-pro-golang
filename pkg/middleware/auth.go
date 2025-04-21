package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/jwt"
)

type key string

const (
	ContextEmailKey  key = "ContextEmailKey"
	ContextUserIDKey key = "ContextUserIDKey"
)

func writeUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func IsAuthed(next http.Handler, jwtService *jwt.JWT) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			writeUnauthorized(w)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		valid, data := jwtService.ParseToken(token)
		if !valid || data == nil {
			writeUnauthorized(w)
			return
		}

		// кладём в контекст и Email, и UserID
		ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)
		ctx = context.WithValue(ctx, ContextUserIDKey, data.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
