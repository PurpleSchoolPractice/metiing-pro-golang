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
	mux.HandleFunc("POST /auth/forgot-password", handler.ForgotPassword())
	mux.HandleFunc("POST /auth/check-token", handler.CheckTokenExpirationDate())
	mux.HandleFunc("PUT /auth/reset-password", handler.ResetPassword())
}

// Register регистрирует пользователя
func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandelBody[RegisterRequest](w, r)
		if err != nil {
			return
		}
		user, err := handler.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		tokenPair, err := handler.AuthService.JWT.GenerateTokenPair(jwt.JWTData{Email: user.Email})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := RegisterResponse{
			UserId:       user.ID,
			Email:        user.Email,
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
			UserId:       jwtData.UserID,
			Email:        jwtData.Email,
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

// ForgotPassword Принимает запрос на восстановление пароля
func (handler *AuthHandler) ForgotPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandelBody[ForgotPasswordRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = handler.AuthService.ForgotPassword(handler.Config, body.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.JsonResponse(w, "A link to reset your password has been sent to your email.", http.StatusOK)
	}
}

// CheckTokenExpirationDate Принимает запрос на проверку срока действия и неиспользованности временного токена для сброса пароля
func (handler *AuthHandler) CheckTokenExpirationDate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandelBody[CheckTokenRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		passwordReset, err := handler.passwordResetRepository.GetActiveToken(body.Token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if passwordReset == nil {
			http.Error(w, "Token not found", http.StatusNotFound)
			return
		}

		res.JsonResponse(w, CheckTokenResponse{
			UserId:    passwordReset.UserID,
			Token:     passwordReset.Token,
			ExpiresAt: passwordReset.ExpiresAt,
		}, http.StatusOK)
	}
}

// ResetPassword Принимает запрос на внесение нового пароля в базу данных
func (handler *AuthHandler) ResetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandelBody[ResetPasswordRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if body.NewPassword == "" || body.Token == "" {
			http.Error(w, "Token and password required", http.StatusBadRequest)
			return
		}

		updatedUser, err := handler.AuthService.ResetPassword(body.Token, body.NewPassword)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.JsonResponse(w, ResetPasswordResponse{
			UserId:    updatedUser.ID,
			UpdatedAt: updatedUser.UpdatedAt,
		}, http.StatusOK)
	}
}
