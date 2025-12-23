package eventParticipant

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/event"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/jwt"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/middleware"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/res"
	"github.com/go-chi/chi/v5"
)

type EventParticipantHandler struct {
	EventParticipantRepository *EventParticipantRepository
	JWTService                 *jwt.JWT
}

type EventParticipantDepsHandler struct {
	EventParticipantRepository *EventParticipantRepository
	JWTService                 *jwt.JWT
}

func NewEventParticipantHandler(mux *chi.Mux, deps EventParticipantDepsHandler) {
	handler := &EventParticipantHandler{
		EventParticipantRepository: deps.EventParticipantRepository,
		JWTService:                 deps.JWTService,
	}
	mux.Handle("POST /event-participant/",
		middleware.IsAuthed(handler.AddEventParticipant(), deps.JWTService))
	mux.Handle("DELETE /event-participant/{id}/event/{event_id}",
		middleware.IsAuthed(handler.DeleteEventParticipant(), deps.JWTService))
	mux.Handle("GET /event-participant/{id}",
		middleware.IsAuthed(handler.GetEventParticipantById(), deps.JWTService))
	mux.Handle("GET /event-participant/user/{user_id}/events",
		middleware.IsAuthed(handler.GetUserEvents(), deps.JWTService))
	mux.Handle("POST /event-participant/is-participant",
		middleware.IsAuthed(handler.IsParticipant(), deps.JWTService))
}

// AddEventParticipant Добавляет нового участника в событие
func (h *EventParticipantHandler) AddEventParticipant() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			EventID uint `json:"event_id"`
			UserID  uint `json:"user_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		userID, err := event.GetUserIDFromContext(r.Context())
		if err != nil {
			switch err {
			case event.EventErrors["missing"]:
				http.Error(w, "user id not found in context", http.StatusUnauthorized)
				return
			case event.EventErrors["type"]:
				http.Error(w, "user id has invalid type", http.StatusInternalServerError)
				return
			default:
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}
		log.Printf(
			"AddEvent: ctxUserID=%d, req.EventID=%d, req.UserID=%d",
			userID, req.EventID, req.UserID,
		)
		isCreator, err := h.EventParticipantRepository.IsEventCreatorById(req.EventID, userID)
		if err != nil {
			http.Error(w, "Error checking if user is creator of event", http.StatusInternalServerError)
			return
		}
		if !isCreator {
			http.Error(w, "User is not creator of event", http.StatusForbidden)
			return
		}

		if err := h.EventParticipantRepository.AddParticipant(req.EventID, req.UserID); err != nil {
			http.Error(w, "Failed to add participant", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// DeleteEventParticipant Удаляет участника из события
func (h *EventParticipantHandler) DeleteEventParticipant() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		participantParam := chi.URLParam(r, "id")
		participantID, err := strconv.ParseUint(participantParam, 10, 64)
		if err != nil {
			http.Error(w, "Invalid participant ID", http.StatusBadRequest)
			return
		}

		eventParam := chi.URLParam(r, "event_id")
		eventID, err := strconv.ParseUint(eventParam, 10, 64)
		if err != nil {
			switch err {
			case event.EventErrors["missing"]:
				http.Error(w, "user id not found in context", http.StatusUnauthorized)
				return
			case event.EventErrors["type"]:
				http.Error(w, "user id has invalid type", http.StatusInternalServerError)
				return
			default:
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}

		userID, err := event.GetUserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, "Error getting user ID from context", http.StatusInternalServerError)
			return
		}
		isCreator, err := h.EventParticipantRepository.IsEventCreatorById(uint(eventID), userID)
		if err != nil {
			http.Error(w, "Error checking if user is creator of event", http.StatusInternalServerError)
			return
		}
		if !isCreator {
			http.Error(w, "User is not creator of event", http.StatusForbidden)
			return
		}

		if err := h.EventParticipantRepository.RemoveParticipant(uint(eventID), uint(participantID)); err != nil {
			http.Error(w, "Failed to remove participant", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// GetEventParticipantById Возвращает список участников события с приглашениями
func (h *EventParticipantHandler) GetEventParticipantById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		participants, err := h.EventParticipantRepository.GetUsersWithInvites(uint(id))
		if err != nil {
			http.Error(w, "Failed to get participants", http.StatusInternalServerError)
			return
		}
		//Собираем ответ
		usersInvite := models.UserStatus{
			UserId: participants.UserID,
			Status: participants.Status,
		}
		res.JsonResponse(w, usersInvite, http.StatusOK)

	}
}

// GetUserEvents Возвращает список событий, в которых участвует пользователь
func (h *EventParticipantHandler) GetUserEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDParam := chi.URLParam(r, "user_id")
		userID, err := strconv.ParseUint(userIDParam, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		events, err := h.EventParticipantRepository.GetUserEvents(uint(userID))
		if err != nil {
			http.Error(w, "Failed to get user events", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(events); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

// IsParticipant Проверяет, является ли пользователь участником события
func (h *EventParticipantHandler) IsParticipant() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			EventID uint `json:"event_id"`
			UserID  uint `json:"user_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		isParticipant, err := h.EventParticipantRepository.IsParticipant(req.EventID, req.UserID)
		if err != nil {
			http.Error(w, "Failed to check participation", http.StatusInternalServerError)
			return
		}

		response := map[string]bool{"is_participant": isParticipant}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
