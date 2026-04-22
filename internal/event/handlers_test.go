package event

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"
	"github.com/go-chi/chi/v5"
)

type MockEventRepository struct {
	FindByIdFunck func(id uint) (*models.Event, error)
}

func (m *MockEventRepository) FindById(id uint) (*models.Event, error) {
	return m.FindByIdFunck(id)
}

func (m *MockEventRepository) Create(event *models.Event) (*models.Event, error) {
	return nil, nil
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
func (m *MockEventRepository) IsUserBusy(userID uint, start time.Time, duration int) bool {
	return false
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
				EventRepository: mockRepo,
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
