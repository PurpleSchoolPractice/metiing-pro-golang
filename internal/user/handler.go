package user

import (
	"net/http"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/convert"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/jwt"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/middleware"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/request"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/res"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type UserHandlerDeps struct {
	UserRepository *UserRepository
	Config         *configs.Config
	JWTService     *jwt.JWT
}
type UserHandler struct {
	UserRepository *UserRepository
	Config         *configs.Config
	JWTService     *jwt.JWT
}

func NewUserHandler(mux *chi.Mux, deps UserHandlerDeps) {
	handler := &UserHandler{
		Config:         deps.Config,
		UserRepository: deps.UserRepository,
		JWTService:     deps.JWTService,
	}
	mux.Handle("GET /users", handler.GetAllUsers())
	mux.Handle("GET /users/{id}", middleware.IsAuthed(handler.GetUserByID(), handler.JWTService))
	mux.Handle("PUT /user/{id}", middleware.IsAuthed(handler.UpdateDataUser(), handler.JWTService))
	mux.Handle("DELETE /user/{id}", middleware.IsAuthed(handler.DeleteUser(), handler.JWTService))
}

// получение всех пользователей
func (handler *UserHandler) GetAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gotAllUsers, err := handler.UserRepository.FindAllUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if gotAllUsers == nil {
			res.JsonResponse(w, "Users not found", http.StatusNotFound)
		}
		res.JsonResponse(w, gotAllUsers, 200)
	}
}

// Получение пользователя по id
func (handler *UserHandler) GetUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := convert.ConvertID(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
		res.JsonResponse(w, user, 200)
	}
}

// Обновление пользователя
func (handler *UserHandler) UpdateDataUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := request.HandelBody[UserUpdateRequest](w, r)
		if err != nil {
			http.Error(w, "Неверный запрос", http.StatusBadRequest)
			return
		}
		userId, err := convert.ConvertID(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		user := &User{
			Model: gorm.Model{
				ID: userId,
			},
			Username: body.Username,
			Password: body.Password,
			Email:    body.Email,
		}
		updatedUser, err := handler.UserRepository.Update(user)
		if err != nil {
			http.Error(w, "We can't to update User %d", int(userId))
			return
		}
		res.JsonResponse(w, updatedUser, 200)
	}
}

// Удаление пользователя
func (handler *UserHandler) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := convert.ConvertID(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		user := &User{
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
