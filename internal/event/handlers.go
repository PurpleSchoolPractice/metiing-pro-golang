package event

import (
	"net/http"
	"strconv"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/eventParticipant"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/convert"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/event"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/jwt"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/middleware"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/request"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/res"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/sendmail"
	"github.com/go-chi/chi/v5"
)

const (
	link string = "http://localhost:8080/event/"
)

type EventHandler struct {
	EventRepository  *EventRepository
	UserRepository   *user.UserRepository
	EventParticipant *eventParticipant.EventParticipantRepository
	JWTService       *jwt.JWT
	Config           *configs.Config
}

type EventHandlerDeps struct {
	EventRepository  *EventRepository
	UserRepository   *user.UserRepository
	EventParticipant *eventParticipant.EventParticipantRepository
	JWTService       *jwt.JWT
	Config           *configs.Config
}

func NewEventHandler(mux *chi.Mux, deps EventHandlerDeps) {
	handler := &EventHandler{
		EventRepository:  deps.EventRepository,
		UserRepository:   deps.UserRepository,
		EventParticipant: deps.EventParticipant,
		JWTService:       deps.JWTService,
		Config:           deps.Config,
	}
	mux.Handle("POST /event/", middleware.IsAuthed(handler.CreateEvent(), handler.JWTService))
	mux.Handle("GET /event/{id}", middleware.IsAuthed(handler.GetEventById(), handler.JWTService))
	mux.Handle("PUT /event/{id}", middleware.IsAuthed(handler.UpdateEvent(), handler.JWTService))
	mux.Handle("DELETE /event/{id}", middleware.IsAuthed(handler.DeleteEvent(), handler.JWTService))
	mux.Handle("GET /event/with-creators", middleware.IsAuthed(handler.GetEventsWithCreators(), handler.JWTService))
	mux.Handle("GET /event/{id}/with-creator", middleware.IsAuthed(handler.GetEventWithCreator(), handler.JWTService))
	mux.Handle("PUT /event/{id}/accept/{userid}", middleware.IsAuthed(handler.Accept(), handler.JWTService))
	mux.Handle("PUT /event/{id}/decline/{userid}", middleware.IsAuthed(handler.Decline(), handler.JWTService))
}

// GetEventById Получает событие по его ID
func (h *EventHandler) GetEventById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id") // получение id из URL
		if idParam == "" {
			http.Error(w, "id not found", http.StatusBadRequest)
			return
		}

		idUint, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			http.Error(w, "Not a valid ID", http.StatusBadRequest)
			return
		}
		id := uint(idUint)

		events, err := h.EventRepository.FindById(id)
		if err != nil {
			http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
			return
		}
		res.JsonResponse(w, events, http.StatusOK)
	}
}

// CreateEvent Создает новое событие
func (h *EventHandler) CreateEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := request.HandelBody[EventRequest](w, r)
		if err != nil {
			http.Error(w, "Неверный запрос", http.StatusBadRequest)
			return
		}
		//валидируем время из запроса
		startTime, err := request.ValidateTime(body.StartDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//создаем новое событие
		newEvent := models.NewEvent(body.Title, body.Description, body.Duration, body.CreatorID, startTime)
		//создаем событие в БД
		createdEvent, err := h.EventRepository.Create(newEvent)
		if err != nil {
			http.Error(w, "Not possible to create new event", http.StatusInternalServerError)
			return
		}
		//логика проверки занятости пользователя
		var userStatusInvate []models.UserStatus
		for _, invUser := range body.InvatedUsers {
			//ищем имя пользователя для ответа по юзер ИД из запроса
			foundUser, err := h.UserRepository.FindByid(invUser.UserId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			//поиск занятости пользователя
			isBusy := h.EventRepository.IsUserBusy(invUser.UserId, startTime, body.Duration)
			//если нашли пересечения
			status := models.StatusBusy
			//если нет пересечений то статус принято и отправляем уведомление на емейл или в лк
			if !isBusy {

				//подготавливаем ссылки
				strEventId := strconv.FormatUint(uint64(createdEvent.ID), 10)
				strUserId := strconv.FormatUint(uint64(invUser.UserId), 10)

				acceptLink := link + strEventId + "/" + "accept" + "/" + strUserId
				declineLink := link + strEventId + "/" + "decline" + "/" + strUserId
				sendmail.SendMail(h.Config, acceptLink, declineLink)

			}
			user := models.UserStatus{
				UserId:   invUser.UserId,
				UserName: foundUser.Username,
				Status:   status,
			}
			userStatusInvate = append(userStatusInvate, user)
		}

		//добавляем участников
		for _, user := range userStatusInvate {
			err := h.EventParticipant.AddParticipant(createdEvent.ID, user.UserId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		//Собираем ответ

		respEvent := &EventResponse{
			Title:       createdEvent.Title,
			Description: createdEvent.Description,
			StartDate:   body.StartDate,
			Duration:    createdEvent.Duration,
			Status:      userStatusInvate,
		}

		res.JsonResponse(w, respEvent, http.StatusCreated)
	}
}

// UpdateEvent Обновляет событие
func (h *EventHandler) UpdateEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := r.Context().Value(middleware.ContextUserIDKey).(uint)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		eventId, err := convert.ParseId(r, "id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hasEvent, err := h.EventRepository.FindById(eventId)
		if err != nil {
			http.Error(w, "Event not found", http.StatusBadRequest)
			return
		}
		//обрабатыввем запрос
		body, err := request.HandelBody[EventRequest](w, r)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		startTime, err := request.ValidateTime(body.StartDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//проверяем является ли юзер создателем события
		if hasEvent.CreatorID != userId {
			http.Error(w, "You are not creator,only creator can update event", http.StatusBadRequest)
			return
		}

		//Заполняем событие новыми данными
		hasEvent.Title = body.Title
		hasEvent.Description = body.Description
		hasEvent.StartDate = startTime
		hasEvent.Duration = body.Duration

		updatedEvent, err := h.EventRepository.Update(hasEvent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//получаем участников события
		partUserEvent, err := h.EventParticipant.GetEventParticipants(updatedEvent.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var userStatusInvate []models.UserStatus
		for _, invUser := range partUserEvent {

			//поиск занятости пользователя
			isBusy := h.EventRepository.IsUserBusy(invUser.ID, startTime, body.Duration)
			//если нашли пересечения
			status := models.StatusBusy
			//если нет пересечений то статус принято
			if !isBusy {

				//подготавливаем ссылки
				strEventId := strconv.FormatUint(uint64(updatedEvent.ID), 10)
				strUserId := strconv.FormatUint(uint64(invUser.ID), 10)

				acceptLink := link + strEventId + "/" + "accept" + "/" + strUserId
				declineLink := link + strEventId + "/" + "decline" + "/" + strUserId
				sendmail.SendMail(h.Config, acceptLink, declineLink)
			}
			user := models.UserStatus{
				UserId:   invUser.ID,
				UserName: invUser.Username,
				Status:   status,
			}
			userStatusInvate = append(userStatusInvate, user)
			//обновляем статусы с участнкиами событий
			h.EventParticipant.UpdateParticipant(&models.EventParticipant{
				EventID: eventId, //берем номер события из пути
				UserID:  user.UserId,
				Status:  status,
			})
		}
		respEvent := &EventResponse{
			Title:       hasEvent.Title,
			Description: hasEvent.Description,
			StartDate:   body.StartDate,
			Duration:    hasEvent.Duration,
			Status:      userStatusInvate,
		}

		res.JsonResponse(w, respEvent, http.StatusOK)
	}
}

// DeleteEvent Удаляет событие
func (h *EventHandler) DeleteEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		if idParam == "" {
			http.Error(w, "id not found", http.StatusBadRequest)
			return
		}

		idUint, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			http.Error(w, "Not a valid ID", http.StatusBadRequest)
			return
		}
		id := uint(idUint)

		_, err = h.EventRepository.FindById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.EventRepository.DeleteById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resDel := &DeleteResponse{
			Delete: true,
		}
		res.JsonResponse(w, resDel, http.StatusOK)
	}
}

// GetEventsWithCreators Получает события вместе с их создателями
func (h *EventHandler) GetEventsWithCreators() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventsWithCreators, err := h.EventRepository.GetEventsWithCreators()
		if eventsWithCreators == nil {
			res.JsonResponse(w, "Not found events", http.StatusOK)
			return
		}
		if err != nil {
			http.Error(w, "Failed to fetch events with creators", http.StatusInternalServerError)
			return
		}

		res.JsonResponse(w, eventsWithCreators, http.StatusOK)
	}
}

// GetEventWithCreator Получает событие вместе с его создателем
func (h *EventHandler) GetEventWithCreator() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventID, err := convert.ParseId(r, "id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userID, err := event.GetUserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, "Error getting user ID from context", http.StatusInternalServerError)
			return
		}

		eventWithCreator, err := h.EventRepository.GetEventWithCreator(eventID, userID)
		if err != nil {
			http.Error(w, "Failed to fetch event with creator", http.StatusInternalServerError)
			return
		}
		res.JsonResponse(w, eventWithCreator, http.StatusOK)

	}
}

func (h *EventHandler) Accept() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//Проверка что только этот пользователь может принять приглашение
		//добавить в список участников с приглашениями
		//прописать json ответ в виде accept

		userId, ok := r.Context().Value(middleware.ContextUserIDKey).(uint)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		eventId, err := convert.ParseId(r, "id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//парсим юзер ИД
		userIdFromUrl, err := convert.ParseId(r, "userid")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//проверяем что юзер является участником события и только он может принять или отклонить
		if userId != userIdFromUrl {
			http.Error(w, "Wrong user", http.StatusConflict)
			return
		}

		//обновляем статус
		updateStatus := &models.EventParticipant{
			EventID: eventId,
			UserID:  userId,
			Status:  models.StatusAccepted,
		}
		updatedStatus, err := h.EventParticipant.UpdateParticipant(updateStatus)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.JsonResponse(w, updatedStatus, http.StatusOK)

	}
}
func (h *EventHandler) Decline() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := r.Context().Value(middleware.ContextUserIDKey).(uint)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		eventId, err := convert.ParseId(r, "id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//парсим юзер ИД
		userIdFromUrl, err := convert.ParseId(r, "userid")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if userId != userIdFromUrl {
			http.Error(w, "Wrong user", http.StatusConflict)
		}

		//обновляем статус
		updateStatus := &models.EventParticipant{
			EventID: eventId,
			UserID:  userId,
			Status:  models.StatusDecline,
		}
		updatedStatus, err := h.EventParticipant.UpdateParticipant(updateStatus)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.JsonResponse(w, updatedStatus, http.StatusOK)

	}
}
