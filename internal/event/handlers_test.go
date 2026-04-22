package event

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"
	"github.com/go-chi/chi/v5"
)

type MockEventRepository struct {
	FindByIdFunck  func(id uint) (*models.Event, error)
	CreateFunc     func(event *models.Event) (*models.Event, error)
	IsUserBusyFunc func(userID uint, start time.Time, duration int) bool
}

func (m *MockEventRepository) FindById(id uint) (*models.Event, error) {
	if m.FindByIdFunck != nil {
		return m.FindByIdFunck(id)
	}
	return nil, nil
}

func (m *MockEventRepository) Create(event *models.Event) (*models.Event, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(event)
	}
	return nil, nil
}

func (m *MockEventRepository) IsUserBusy(userID uint, start time.Time, duration int) bool {
	if m.IsUserBusyFunc != nil {
		return m.IsUserBusyFunc(userID, start, duration)
	}
	return false
}

func (m *MockEventRepository) Update(event *models.Event) (*models.Event, error) {
	return nil, nil
}

func (m *MockEventRepository) DeleteById(id uint) error {
	return nil
}
func (m *MockEventRepository) GetEventsWithCreators() ([]models.Event, error) {
	return nil, nil
}

func (m *MockEventRepository) GetEventWithCreator(eventID, userID uint) (*models.Event, error) {
	return nil, nil
}

func (m *MockEventRepository) FindAllByCreatorId(id uint) ([]models.Event, error) {
	return nil, nil
}

type MockUserRepository struct {
	FindByIdFunc func(id uint) (*models.User, error)
}

func (m *MockUserRepository) FindById(id uint) (*models.User, error) {
	if m.FindByIdFunc != nil {
		return m.FindByIdFunc(id)
	}
	return nil, nil
}

func (m *MockUserRepository) Create(user *models.User) (*models.User, error) {
	return nil, nil
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	return nil, nil
}

func (m *MockUserRepository) FindAllUsers(limit, offset int, search string) ([]models.UserResponse, int64, error) {
	return nil, 0, nil
}

func (m *MockUserRepository) Update(user *models.User) (*models.User, error) {
	return nil, nil
}

func (m *MockUserRepository) DeleteById(id uint) error {
	return nil
}

type MockEventParticipantRepository struct {
	AddParticipantFunc    func(eventID, userID uint) error
	UpdateParticipantFunc func(participant *models.EventParticipant) (*models.EventParticipant, error)
}

func (m *MockEventParticipantRepository) AddParticipant(eventID, userID uint) error {
	if m.AddParticipantFunc != nil {
		return m.AddParticipantFunc(eventID, userID)
	}
	return nil
}

func (m *MockEventParticipantRepository) RemoveParticipant(eventID, userID uint) error {
	return nil
}

func (m *MockEventParticipantRepository) GetEventParticipants(eventID uint) ([]models.User, error) {
	return nil, nil
}

func (m *MockEventParticipantRepository) GetUserEvents(userID uint) ([]models.Event, error) {
	return nil, nil
}

func (m *MockEventParticipantRepository) IsParticipant(eventID, userID uint) (bool, error) {
	return false, nil
}
func (m *MockEventParticipantRepository) UpdateParticipant(participant *models.EventParticipant) (*models.EventParticipant, error) {
	if m.UpdateParticipantFunc != nil {
		return m.UpdateParticipantFunc(participant)
	}
	return participant, nil
}
func TestGetEventByID(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockReturn     *models.Event
		mockError      error
		expectedStatus int
	}{
		{
			name:           "id отсутствует",
			id:             "",
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "id не число",
			id:             "abc",
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "repo error",
			id:             "1",
			mockReturn:     nil,
			mockError:      errors.New("ошибка БД"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "успех",
			id:   "1",
			mockReturn: &models.Event{
				Title:       "Тест",
				Description: "Описание",
				Duration:    60,
				CreatorID:   100,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &MockEventRepository{
				FindByIdFunck: func(id uint) (*models.Event, error) {
					return tt.mockReturn, tt.mockError
				},
			}

			hendler := &EventHandler{
				EventRepository:  mockRepo,
				UserRepository:   &MockUserRepository{},
				EventParticipant: &MockEventParticipantRepository{},
			}

			req := httptest.NewRequest("GET", "/event/"+tt.id, nil)

			rctx := chi.NewRouteContext()
			if tt.id != "" {
				rctx.URLParams.Add("id", tt.id)
			}

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()

			hendler.GetEventById()(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Ожидался статус %d, получен %d", tt.expectedStatus, status)
			}

			if tt.name == "успех" {
				var result models.Event

				if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
					t.Fatalf("Не удалось распарсить JSON: %v", err)
				}
				if result.ID != tt.mockReturn.ID {
					t.Errorf("Ожидался ID %d, получен %d", tt.mockReturn.ID, result.ID)
				}
			}

		})
	}

}

func TestCreateEventHandler(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        string
		mockCreate         func(event *models.Event) (*models.Event, error)
		mockFindById       func(id uint) (*models.User, error)
		mockAddParticipant func(uint, uint) error
		expectedStatus     int
	}{
		{
			name:           "битый JSON",
			requestBody:    `{invalid json}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "невалидный start_date",
			requestBody: `{
				"title": "Test",
				"start_date": "invalid-date",
				"duration": 60,
				"creator_id": 1 }`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "ошибка создания в репо",
			requestBody: `{
				"title": "Test",
				"start_date": "2025-05-09 11:00",
				"duration": 60,
				"creator_id": 1}`,
			mockCreate: func(event *models.Event) (*models.Event, error) {
				return nil, errors.New("ошибка бд")
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "ошибка FindById пользователя",
			requestBody: `{
				"title": "Test",
				"start_date": "2025-05-09 11:00",
				"duration": 60,
				"creator_id": 1,
				"invated_users": [{"user_id": 10}]}`,
			mockCreate: func(event *models.Event) (*models.Event, error) {
				return &models.Event{CreatorID: 1}, nil
			},
			mockFindById: func(id uint) (*models.User, error) {
				return nil, errors.New("пользователь не найден")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "ошибка AddParticipant",
			requestBody: `{
				"title": "Test",
				"start_date": "2025-05-09 11:00",
				"duration": 60,
				"creator_id": 1,
				"invated_users": [{"user_id": 10}]}`,
			mockCreate: func(event *models.Event) (*models.Event, error) {
				return &models.Event{CreatorID: 1}, nil
			},
			mockFindById: func(id uint) (*models.User, error) {
				user := &models.User{
					Username: "test",
					Email:    "test@example.com",
					Password: "password",
				}
				user.ID = 10

				return user, nil
			},
			mockAddParticipant: func(eventID, userID uint) error {
				return errors.New("ошибка добавления")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "успех",
			requestBody: `{
				"title": "Test Event",
				"description": "Test Description",
				"start_date": "2025-05-09 11:00",
				"duration": 60,
				"creator_id": 1,
				"invated_users": [{"user_id": 10}]}`,
			mockCreate: func(event *models.Event) (*models.Event, error) {
				return &models.Event{
					CreatorID:   1,
					Title:       "Test Event",
					Description: "Test Description",
					Duration:    60,
				}, nil
			},
			mockFindById: func(id uint) (*models.User, error) {
				user := &models.User{
					Username: "test",
					Email:    "test@example.com",
					Password: "password",
				}
				user.ID = 10

				return user, nil
			},
			mockAddParticipant: func(eventID, userID uint) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &EventHandler{
				EventRepository: &MockEventRepository{
					CreateFunc: tt.mockCreate,
				},
				UserRepository: &MockUserRepository{
					FindByIdFunc: tt.mockFindById,
				},
				EventParticipant: &MockEventParticipantRepository{
					AddParticipantFunc: tt.mockAddParticipant,
				},
				Config: &configs.Config{},
			}

			req := httptest.NewRequest("POST", "/event", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler.CreateEvent()(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Ожидался статус %d, получен %d", tt.expectedStatus, rr.Code)
			}

			if tt.name == "успех" {
				var result EventResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
					t.Fatalf("Не удалось распарсить JSON: %v", err)
				}

				if result.Title != "Test Event" {
					t.Errorf("Ожидался title 'Test Event', получен %q", result.Title)
				}

				if result.Description != "Test Description" {
					t.Errorf("Ожидался description 'Test Description' , получен %q", result.Description)
				}

				if result.Duration != 60 {
					t.Errorf("Ожидался  duration 60, получен %d", result.Duration)
				}

				if len(result.Status) != 1 {
					t.Errorf("Ожидался 1 статус, получено %d", len(result.Status))
				}

			}

		})
	}
}
