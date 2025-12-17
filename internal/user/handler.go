package user

import (
	"net/http"
	"strconv"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/convert"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/jwt"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/middleware"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/request"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/res"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandlerDeps struct {
	UserRepository *UserRepository
	JWTService     *jwt.JWT
}
type UserHandler struct {
	UserRepository *UserRepository
	JWTService     *jwt.JWT
}

func NewUserHandler(mux *chi.Mux, deps UserHandlerDeps) {
	handler := &UserHandler{
		UserRepository: deps.UserRepository,
		JWTService:     deps.JWTService,
	}
	mux.Handle("GET /users", handler.GetAllUsers())
	mux.Handle("GET /user/{id}", middleware.IsAuthed(handler.GetUserByID(), handler.JWTService))
	mux.Handle("PUT /user/{id}", middleware.IsAuthed(handler.UpdateDataUser(), handler.JWTService))
	mux.Handle("DELETE /user/{id}", middleware.IsAuthed(handler.DeleteUser(), handler.JWTService))
}

// получение всех пользователей
func (handler *UserHandler) GetAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		limit := 20
		offset := 0

		limitStr := r.URL.Query().Get("limit")
		if limitStr != "" {
			limitInt, err := strconv.Atoi(limitStr) 
			if err != nil {
				http.Error(w, "Invalid limit param", http.StatusBadRequest)
				return
			} 
			limit = limitInt
		}

		offsetStr := r.URL.Query().Get("offset")
		if offsetStr != "" {
			offsetInt, err := strconv.Atoi(offsetStr) 
			if err != nil {
				http.Error(w, "Invalid offset param", http.StatusBadRequest)
				return
			} 
			offset = offsetInt
		}

		search := r.URL.Query().Get("q")

		if limit < 1 || limit > 100 {
			http.Error(w, "Invalid limit param, max 100 value", http.StatusBadRequest)
			return
		}

		if offset < 0 {
			http.Error(w, "Invalid offset param, min 0 value", http.StatusBadRequest)
			return
		}

		if len(search) > 200 {
			http.Error(w, "Invalid search param, max 200 characters", http.StatusBadRequest)
			return
		}

		allUsers, total, err := handler.UserRepository.FindAllUsers(limit, offset, search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if allUsers == nil {
			res.JsonResponse(w, "Users not found", http.StatusNotFound)
			return
		}

		res.JsonResponse(w, UserPaginatedResponse{
			Items: allUsers,
			Total: total,
			Limit: limit,
			Offset: offset,
		}, 200)
	}
}

// Получение пользователя по id
func (handler *UserHandler) GetUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//Парсим ID из строки
		id, err := convert.ParseId(r, "id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := handler.UserRepository.FindByid(id)
		if err != nil {
			http.Error(w, "Failed search user by id", http.StatusBadRequest)
			return
		}
		if user == nil {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}

		userResponse := models.ToUserResponse(user)

		res.JsonResponse(w, userResponse, 200)
	}
}

// Обновление пользователя
func (handler *UserHandler) UpdateDataUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIdContext, ok := r.Context().Value(middleware.ContextUserIDKey).(uint)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		body, err := request.HandelBody[UserUpdateRequest](w, r)
		if err != nil {
			http.Error(w, "Неверный запрос", http.StatusBadRequest)
			return
		}
		//Парсим ID из строки
		userId, err := convert.ParseId(r, "id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}
		//Заполняем данные юзера для обновления
		user := &models.User{
			Model: gorm.Model{
				ID: userId,
			},
			Username: body.Username,
			Password: string(hashedPassword),
			Email:    body.Email,
		}
		var updatedUser *models.User
		//Проверяем что обновляем именно авторизованного юзера, а не кого другого
		if userIdContext == userId {
			updatedUser, err = handler.UserRepository.Update(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "Can not Update. Different user", http.StatusBadRequest)
			return
		}

		userResponse := models.ToUserResponse(updatedUser)

		res.JsonResponse(w, userResponse, 200)
	}
}

// Удаление пользователя
func (handler *UserHandler) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := convert.ParseId(r, "id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		user := &models.User{
			Model: gorm.Model{
				ID: id,
			},
		}
		err = handler.UserRepository.Delete(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		res.JsonResponse(w, "User deleted", 200)
	}
}
