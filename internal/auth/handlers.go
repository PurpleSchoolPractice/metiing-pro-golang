package auth

import (
	"net/http"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/jwt"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/request"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/res"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	*configs.Config
	*AuthService
}

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

// NewAuthHandler инициализирует хендлеры для аутентификации
func NewAuthHandler(mux *chi.Mux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	mux.HandleFunc("POST /auth/register", handler.Register())
	mux.HandleFunc("POST /auth/login", handler.Login())
	mux.HandleFunc("POST /auth/refresh", handler.RefreshToken())
}

// Register регистрирует пользователя
func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandelBody[RegisterRequest](w, r)
		if err != nil {
			return
		}
		email, err := handler.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		tokenPair, err := handler.AuthService.JWT.GenerateTokenPair(jwt.JWTData{Email: email})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := RegisterResponse{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
		}
		res.JsonResponse(w, data, http.StatusCreated)
	}
}

// Login авторизует пользователя
func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandelBody[LoginRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		jwtData, err := handler.AuthService.Login(body.Email, body.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		tokenPair, err := handler.AuthService.JWT.GenerateTokenPair(jwtData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.JsonResponse(w, LoginResponse{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
		}, http.StatusOK)
	}
}

func (handler *AuthHandler) RefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandelBody[RefreshTokenRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tokenPair, err := handler.AuthService.RefreshTokens(body.RefreshToken, body.AccessToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		data := RefreshTokenResponse{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
		}
		res.JsonResponse(w, data, http.StatusOK)
	}
}
