package event

import (
	"net/http"
	"strconv"
	"time"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/convert"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/jwt"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/middleware"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/request"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/res"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type EventHandler struct {
	EventRepository *EventRepository
	JWTService      *jwt.JWT
}

type EventHandlerDeps struct {
	EventRepository *EventRepository
	JWTService      *jwt.JWT
}

func NewEventHandler(mux *chi.Mux, deps EventHandlerDeps) {
	handler := &EventHandler{
		EventRepository: deps.EventRepository,
		JWTService:      deps.JWTService,
	}
	mux.Handle("POST /event/", middleware.IsAuthed(handler.CreateEvent(), handler.JWTService))
	mux.Handle("GET /event/{id}", middleware.IsAuthed(handler.GetEventById(), handler.JWTService))
	mux.Handle("PUT /event/{id}", middleware.IsAuthed(handler.UpdateEvent(), handler.JWTService))
	mux.Handle("DELETE /event/{id}", middleware.IsAuthed(handler.DeleteEvent(), handler.JWTService))
	mux.Handle("GET /event/with-creators", middleware.IsAuthed(handler.GetEventsWithCreators(), handler.JWTService))
	mux.Handle("GET /event/{id}/with-creator", middleware.IsAuthed(handler.GetEventWithCreator(), handler.JWTService))
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

		newEvent := NewEvent(body.Title, body.Description, body.CreatorID, body.EventDate)
		newEvent.OwnerID = body.CreatorID

		createdEvent, err := h.EventRepository.Create(newEvent)
		if err != nil {
			http.Error(w, "Not possible to create new event", http.StatusInternalServerError)
			return
		}

		res.JsonResponse(w, createdEvent, http.StatusCreated)
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

		eventId, err := convert.ParseId(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		hasEventId, err := h.EventRepository.FindById(eventId)
		if err != nil {
			http.Error(w, "Event not found", http.StatusBadRequest)
			return
		}
		body, err := request.HandelBody[EventRequest](w, r)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		newEvent := &Event{
			Model:       gorm.Model{ID: hasEventId.ID},
			Title:       body.Title,
			Description: body.Description,
			EventDate:   time.Now(),
			CreatorID:   userId,
			OwnerID:     userId,
		}

		updatedEvent, err := h.EventRepository.Update(newEvent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.JsonResponse(w, updatedEvent, http.StatusCreated)
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

		eventWithCreator, err := h.EventRepository.GetEventWithCreator(id)
		if err != nil {
			http.Error(w, "Failed to fetch event with creator", http.StatusInternalServerError)
		}
		res.JsonResponse(w, eventWithCreator, http.StatusOK)

	}
}
