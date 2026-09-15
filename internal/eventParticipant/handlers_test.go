package eventParticipant

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/middleware"
	"github.com/go-chi/chi/v5"
)

type MockEventParticipantRepository struct {
	AddParticipantFunc       func(eventID, userID uint) error
	IsEventCreatorByIdFunc   func(eventID, userID uint) (bool, error)
	RemoveParticipantFunc    func(eventID, userID uint) error
	GetUsersWithInvitesFunc  func(eventID uint) ([]models.EventParticipant, error)
	GetEventParticipantsFunc func(eventID uint) ([]models.User, error)
	GetUserEventsFunc        func(userID uint) ([]models.Event, error)
	IsParticipantFunc        func(eventID, userID uint) (bool, error)
	UpdateParticipantFunc    func(participant *models.EventParticipant) (*models.EventParticipant, error)
}

func (m *MockEventParticipantRepository) IsEventCreatorById(eventID, userID uint) (bool, error) {
	if m.IsEventCreatorByIdFunc != nil {
		return m.IsEventCreatorByIdFunc(eventID, userID)
	}
	return false, nil
}

func (m *MockEventParticipantRepository) AddParticipant(eventID, userID uint) error {
	if m.AddParticipantFunc != nil {
		return m.AddParticipantFunc(eventID, userID)
	}
	return nil
}

func (m *MockEventParticipantRepository) RemoveParticipant(eventID, userID uint) error {
	if m.RemoveParticipantFunc != nil {
		return m.RemoveParticipantFunc(eventID, userID)
	}
	return nil
}

func (m *MockEventParticipantRepository) GetUsersWithInvites(eventID uint) ([]models.EventParticipant, error) {
	if m.GetUsersWithInvitesFunc != nil {
		return m.GetUsersWithInvitesFunc(eventID)
	}
	return nil, nil
}

func (m *MockEventParticipantRepository) GetEventParticipants(eventID uint) ([]models.User, error) {
	if m.GetEventParticipantsFunc != nil {
		return m.GetEventParticipantsFunc(eventID)
	}
	return nil, nil
}

func (m *MockEventParticipantRepository) GetUserEvents(userID uint) ([]models.Event, error) {
	if m.GetUserEventsFunc != nil {
		return m.GetUserEventsFunc(userID)
	}
	return nil, nil
}

func (m *MockEventParticipantRepository) IsParticipant(eventID, userID uint) (bool, error) {
	if m.IsParticipantFunc != nil {
		return m.IsParticipantFunc(eventID, userID)
	}
	return false, nil
}

func (m *MockEventParticipantRepository) UpdateParticipant(
	participant *models.EventParticipant,
) (*models.EventParticipant, error) {
	if m.UpdateParticipantFunc != nil {
		return m.UpdateParticipantFunc(participant)
	}
	return nil, nil
}

func TestAddEventParticipantInvalidJSON(t *testing.T) {
	mock := &MockEventParticipantRepository{}

	handler := &EventParticipantHandler{
		EventParticipantRepository: mock,
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/event-participant/",
		strings.NewReader("{bad json"),
	)

	rec := httptest.NewRecorder()

	h := handler.AddEventParticipant()
	h(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestAddEventParticipantMissingUser(t *testing.T) {
	mock := &MockEventParticipantRepository{}

	handler := &EventParticipantHandler{
		EventParticipantRepository: mock,
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/event-participant/",
		strings.NewReader(`{"event_id":1,"user_id":2}`),
	)

	rec := httptest.NewRecorder()

	h := handler.AddEventParticipant()
	h(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}

func TestAddEventParticipantCreatorCheckError(t *testing.T) {
	mock := &MockEventParticipantRepository{
		IsEventCreatorByIdFunc: func(eventID, userID uint) (bool, error) {
			return false, errors.New("repository error")
		},
	}

	handler := &EventParticipantHandler{
		EventParticipantRepository: mock,
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/event-participant/",
		strings.NewReader(`{"event_id":1,"user_id":2}`),
	)

	ctx := context.WithValue(
		req.Context(),
		middleware.ContextUserIDKey,
		uint(1),
	)

	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	h := handler.AddEventParticipant()
	h(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}
}

func TestAddEventParticipantForbidden(t *testing.T) {
	mock := &MockEventParticipantRepository{
		IsEventCreatorByIdFunc: func(eventID, userID uint) (bool, error) {
			return false, nil
		},
	}

	handler := &EventParticipantHandler{
		EventParticipantRepository: mock,
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/event-participant/",
		strings.NewReader(`{"event_id":1,"user_id":2}`),
	)

	ctx := context.WithValue(
		req.Context(),
		middleware.ContextUserIDKey,
		uint(1),
	)

	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	h := handler.AddEventParticipant()
	h(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusForbidden,
			rec.Code,
		)
	}
}

func TestAddEventParticipantRepositoryError(t *testing.T) {
	mock := &MockEventParticipantRepository{
		IsEventCreatorByIdFunc: func(eventID, userID uint) (bool, error) {
			return true, nil
		},
		AddParticipantFunc: func(eventID, userID uint) error {
			return errors.New("repository error")
		},
	}

	handler := &EventParticipantHandler{
		EventParticipantRepository: mock,
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/event-participant/",
		strings.NewReader(`{"event_id":1,"user_id":2}`),
	)

	ctx := context.WithValue(
		req.Context(),
		middleware.ContextUserIDKey,
		uint(1),
	)

	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	h := handler.AddEventParticipant()
	h(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}
}

func TestAddEventParticipantSuccess(t *testing.T) {
	mock := &MockEventParticipantRepository{
		IsEventCreatorByIdFunc: func(eventID, userID uint) (bool, error) {
			return true, nil
		},
		AddParticipantFunc: func(eventID, userID uint) error {
			return nil
		},
	}

	handler := &EventParticipantHandler{
		EventParticipantRepository: mock,
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/event-participant/",
		strings.NewReader(`{"event_id":1,"user_id":2}`),
	)

	ctx := context.WithValue(
		req.Context(),
		middleware.ContextUserIDKey,
		uint(1),
	)

	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	h := handler.AddEventParticipant()
	h(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			rec.Code,
		)
	}
}

func TestDeleteEventParticipantInvalidParticipantID(t *testing.T) {
	handler := &EventParticipantHandler{
		EventParticipantRepository: &MockEventParticipantRepository{},
	}

	req := httptest.NewRequest(
		http.MethodDelete,
		"/event-participant/bad/event/1",
		nil,
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "bad")
	rctx.URLParams.Add("event_id", "1")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	rec := httptest.NewRecorder()

	h := handler.DeleteEventParticipant()
	h(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestDeleteEventParticipantInvalidEventID(t *testing.T) {
	handler := &EventParticipantHandler{
		EventParticipantRepository: &MockEventParticipantRepository{},
	}

	req := httptest.NewRequest(
		http.MethodDelete,
		"/event-participant/2/event/bad",
		nil,
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "2")
	rctx.URLParams.Add("event_id", "bad")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	rec := httptest.NewRecorder()

	h := handler.DeleteEventParticipant()
	h(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestDeleteEventParticipantMissingUser(t *testing.T) {
	handler := &EventParticipantHandler{
		EventParticipantRepository: &MockEventParticipantRepository{},
	}

	req := httptest.NewRequest(
		http.MethodDelete,
		"/event-participant/2/event/1",
		nil,
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "2")
	rctx.URLParams.Add("event_id", "1")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	rec := httptest.NewRecorder()

	h := handler.DeleteEventParticipant()
	h(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}

func TestDeleteEventParticipantForbidden(t *testing.T) {
	mock := &MockEventParticipantRepository{
		IsEventCreatorByIdFunc: func(eventID, userID uint) (bool, error) {
			return false, nil
		},
	}

	handler := &EventParticipantHandler{
		EventParticipantRepository: mock,
	}

	req := httptest.NewRequest(
		http.MethodDelete,
		"/event-participant/2/event/1",
		nil,
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "2")
	rctx.URLParams.Add("event_id", "1")

	ctx := context.WithValue(req.Context(), middleware.ContextUserIDKey, uint(1))
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	h := handler.DeleteEventParticipant()
	h(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusForbidden,
			rec.Code,
		)
	}
}

func TestDeleteEventParticipantCreatorCheckError(t *testing.T) {
	mock := &MockEventParticipantRepository{
		IsEventCreatorByIdFunc: func(eventID, userID uint) (bool, error) {
			return false, errors.New("repository error")
		},
	}

	handler := &EventParticipantHandler{
		EventParticipantRepository: mock,
	}

	req := httptest.NewRequest(
		http.MethodDelete,
		"/event-participant/2/event/1",
		nil,
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "2")
	rctx.URLParams.Add("event_id", "1")

	ctx := context.WithValue(req.Context(), middleware.ContextUserIDKey, uint(1))
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	h := handler.DeleteEventParticipant()
	h(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}
}

func TestDeleteEventParticipantRepositoryError(t *testing.T) {
	mock := &MockEventParticipantRepository{
		IsEventCreatorByIdFunc: func(eventID, userID uint) (bool, error) {
			return true, nil
		},
		RemoveParticipantFunc: func(eventID, userID uint) error {
			return errors.New("repository error")
		},
	}

	handler := &EventParticipantHandler{
		EventParticipantRepository: mock,
	}

	req := httptest.NewRequest(
		http.MethodDelete,
		"/event-participant/2/event/1",
		nil,
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "2")
	rctx.URLParams.Add("event_id", "1")

	ctx := context.WithValue(req.Context(), middleware.ContextUserIDKey, uint(1))
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	h := handler.DeleteEventParticipant()
	h(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}
}

func TestDeleteEventParticipantSuccess(t *testing.T) {
	mock := &MockEventParticipantRepository{
		IsEventCreatorByIdFunc: func(eventID, userID uint) (bool, error) {
			return true, nil
		},
		RemoveParticipantFunc: func(eventID, userID uint) error {
			return nil
		},
	}

	handler := &EventParticipantHandler{
		EventParticipantRepository: mock,
	}

	req := httptest.NewRequest(
		http.MethodDelete,
		"/event-participant/2/event/1",
		nil,
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "2")
	rctx.URLParams.Add("event_id", "1")

	ctx := context.WithValue(req.Context(), middleware.ContextUserIDKey, uint(1))
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	h := handler.DeleteEventParticipant()
	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}
}

func TestGetEventParticipantByIdInvalidID(t *testing.T) {
	handler := &EventParticipantHandler{
		EventParticipantRepository: &MockEventParticipantRepository{},
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/event-participant/bad",
		nil,
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "bad")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	rec := httptest.NewRecorder()

	h := handler.GetEventParticipantById()
	h(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestGetEventParticipantByIdRepositoryError(t *testing.T) {
	mock := &MockEventParticipantRepository{
		GetUsersWithInvitesFunc: func(eventID uint) ([]models.EventParticipant, error) {
			return nil, errors.New("repository error")
		},
	}

	handler := &EventParticipantHandler{
		EventParticipantRepository: mock,
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/event-participant/1",
		nil,
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	rec := httptest.NewRecorder()

	h := handler.GetEventParticipantById()
	h(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}
}

func TestGetEventParticipantByIdSuccess(t *testing.T) {
	mock := &MockEventParticipantRepository{
		GetUsersWithInvitesFunc: func(eventID uint) ([]models.EventParticipant, error) {
			return []models.EventParticipant{
				{
					UserID: 2,
				},
				{
					UserID: 3,
				},
			}, nil
		},
	}

	handler := &EventParticipantHandler{
		EventParticipantRepository: mock,
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/event-participant/1",
		nil,
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	rec := httptest.NewRecorder()

	h := handler.GetEventParticipantById()
	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	var response []models.UserStatus

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 2 {
		t.Fatalf(
			"expected 2 participants, got %d",
			len(response),
		)
	}

	if response[0].UserId != 2 {
		t.Fatalf(
			"expected first user id %d, got %d",
			2,
			response[0].UserId,
		)
	}

	if response[1].UserId != 3 {
		t.Fatalf(
			"expected second user id %d, got %d",
			3,
			response[1].UserId,
		)
	}
}

func TestGetUserEventsInvalidUserID(t *testing.T) {
	handler := &EventParticipantHandler{
		EventParticipantRepository: &MockEventParticipantRepository{},
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/event-participant/user/bad/events",
		nil,
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("user_id", "bad")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	rec := httptest.NewRecorder()

	h := handler.GetUserEvents()
	h(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestGetUserEventsRepositoryError(t *testing.T) {
	mock := &MockEventParticipantRepository{
		GetUserEventsFunc: func(userID uint) ([]models.Event, error) {
			return nil, errors.New("repository error")
		},
	}

	handler := &EventParticipantHandler{
		EventParticipantRepository: mock,
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/event-participant/user/1/events",
		nil,
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("user_id", "1")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	rec := httptest.NewRecorder()

	h := handler.GetUserEvents()
	h(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}
}

func TestGetUserEventsSuccess(t *testing.T) {
	mock := &MockEventParticipantRepository{
		GetUserEventsFunc: func(userID uint) ([]models.Event, error) {
			return []models.Event{
				{},
				{},
			}, nil
		},
	}

	handler := &EventParticipantHandler{
		EventParticipantRepository: mock,
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/event-participant/user/1/events",
		nil,
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("user_id", "1")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	rec := httptest.NewRecorder()

	h := handler.GetUserEvents()
	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	var response []models.Event

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 2 {
		t.Fatalf(
			"expected 2 events, got %d",
			len(response),
		)
	}
}

func TestIsParticipantInvalidJSON(t *testing.T) {
	handler := &EventParticipantHandler{
		EventParticipantRepository: &MockEventParticipantRepository{},
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/event-participant/is-participant",
		strings.NewReader("{bad json"),
	)

	rec := httptest.NewRecorder()

	h := handler.IsParticipant()
	h(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestIsParticipantRepositoryError(t *testing.T) {
	mock := &MockEventParticipantRepository{
		IsParticipantFunc: func(eventID, userID uint) (bool, error) {
			return false, errors.New("repository error")
		},
	}

	handler := &EventParticipantHandler{
		EventParticipantRepository: mock,
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/event-participant/is-participant",
		strings.NewReader(`{"event_id":1,"user_id":2}`),
	)

	rec := httptest.NewRecorder()

	h := handler.IsParticipant()
	h(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}
}

func TestIsParticipantSuccess(t *testing.T) {
	mock := &MockEventParticipantRepository{
		IsParticipantFunc: func(eventID, userID uint) (bool, error) {
			return true, nil
		},
	}

	handler := &EventParticipantHandler{
		EventParticipantRepository: mock,
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/event-participant/is-participant",
		strings.NewReader(`{"event_id":1,"user_id":2}`),
	)

	rec := httptest.NewRecorder()

	h := handler.IsParticipant()
	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	var response map[string]bool

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["is_participant"] != true {
		t.Fatalf(
			"expected is_participant true, got %v",
			response["is_participant"],
		)
	}
}
