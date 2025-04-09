package github

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func setupMockServer(resp string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(resp))
	}))
}

func TestGetUserEvents(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		mockResponse   string
		mockStatusCode int
		expectedEvents []Event
		wantErr        bool
	}{
		{
			name:           "successful response",
			username:       "testuser",
			mockResponse:   `[{"type":"PushEvent","repo":{"name":"testuser/testrepo"},"payload":{}}]`,
			mockStatusCode: 200,
			expectedEvents: []Event{
				{
					Type: "PushEvent",
					Repo: EventRepo{
						Name: "testuser/testrepo",
					},
					Payload: map[string]interface{}{},
				},
			},
		},
		{
			name:           "not found",
			username:       "nonexistentuser",
			mockResponse:   `{"message":"Not Found"}`,
			mockStatusCode: 404,
			expectedEvents: nil,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := setupMockServer(tt.mockResponse, tt.mockStatusCode)
			defer mockServer.Close()

			httpClient := &http.Client{}

			service := NewService(httpClient)
			service.baseURL = mockServer.URL

			events, err := service.GetUserEvents(tt.username)

			if (err != nil) && !tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(events, tt.expectedEvents) {
				t.Errorf("expected events: %v, got: %v", tt.expectedEvents, events)
			}
		})
	}
}
