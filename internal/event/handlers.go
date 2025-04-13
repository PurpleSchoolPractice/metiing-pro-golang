package event

import (
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/jwt"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/middleware"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/request"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/res"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
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
	handler := EventHandler{
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
		_, ok := r.Context().Value(middleware.ContextEmailKey).(string)
		if !ok {
			http.Error(w, "Not possible to create new event", http.StatusInternalServerError)
		}
		body, err := request.HandelBody[EventRequest](w, r)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
		newEvent := NewEvent(body.Title, body.Description, body.CreatorID, body.EventDate)

		updatedEvent, err := h.EventRepository.Update(newEvent)
		if err != nil {
			http.Error(w, "Not possible to update new event", http.StatusInternalServerError)
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
			http.Error(w, "Not possible to delete event", http.StatusInternalServerError)
			return
		}
		err = h.EventRepository.DeleteById(id)
		if err != nil {
			http.Error(w, "Not possible to delete event", http.StatusInternalServerError)
			return
		}
		res.JsonResponse(w, true, http.StatusOK)
	}
}

// GetEventsWithCreators Получает события вместе с их создателями
func (h *EventHandler) GetEventsWithCreators() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventsWithCreators, err := h.EventRepository.GetEventsWithCreators()
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
