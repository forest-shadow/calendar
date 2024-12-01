package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/forest-shadow/calendar/internal/events"
	"github.com/forest-shadow/calendar/test/mocks"
)

type testCase struct {
	name           string
	event          events.Event
	buildStubs     func(service *mocks.MockeventService)
	checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	expectedStatus int
}

func TestCreateEvent(t *testing.T) {
	testID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")
	startTime, _ := time.Parse(time.RFC3339, "2024-01-01T10:00:00Z")
	endTime, _ := time.Parse(time.RFC3339, "2024-01-01T11:00:00Z")
	notifyBefore := "PT10S"
	invalidNotifyBefore := "123"
	description := "Test Description"

	validEvent := events.Event{
		ID:           testID,
		Title:        "Test Event",
		StartTime:    startTime,
		EndTime:      endTime,
		Description:  &description,
		UserID:       userID,
		NotifyBefore: &notifyBefore,
	}

	testCases := []testCase{
		{
			name:  "Success",
			event: validEvent,
			buildStubs: func(service *mocks.MockeventService) {
				service.EXPECT().
					Create(gomock.Any(), validEvent).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)
				assert.Equal(t, "", recorder.Body.String())
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid NotifyBefore Format",
			event: events.Event{
				ID:           testID,
				Title:        "Test Event",
				StartTime:    startTime,
				EndTime:      endTime,
				Description:  &description,
				UserID:       userID,
				NotifyBefore: &invalidNotifyBefore,
			},
			buildStubs: func(service *mocks.MockeventService) {
				// No expectations as it should fail before service call
				service.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockeventService(ctrl)
			tc.buildStubs(service)

			// Create router and register handler
			router := NewRouter(zap.NewNop().Sugar(), service)

			// Create request
			body, err := json.Marshal(tc.event)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/events", bytes.NewBuffer(body))
			require.NoError(t, err)

			recorder := httptest.NewRecorder()

			// Use the router to serve the request
			router.ServeHTTP(recorder, req)

			// Check response
			tc.checkResponse(t, recorder)
		})
	}
}
