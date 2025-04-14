package github

import (
	"fmt"
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
		{
			name:           "internal server error",
			username:       "testuser",
			mockResponse:   `{"message":"Internal Server Error"}`,
			mockStatusCode: 500,
			expectedEvents: nil,
			wantErr:        true,
		},
		{
			name:           "invalid json response",
			username:       "testuser",
			mockResponse:   `{"invalid_json":`,
			mockStatusCode: 200,
			expectedEvents: nil,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := setupMockServer(tt.mockResponse, tt.mockStatusCode)
			defer mockServer.Close()

			service := NewService(WithBaseURL(mockServer.URL))

			events, err := service.getUserEvents(tt.username)

			if (err != nil) && !tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(events, tt.expectedEvents) {
				t.Errorf("expected events: %v, got: %v", tt.expectedEvents, events)
			}
		})
	}
}

func TestPushEvent_Summarize(t *testing.T) {
	tests := []struct {
		name     string
		event    Event
		expected string
	}{
		{
			name: "valid push event",
			event: Event{
				Type: PushEvent,
				Repo: EventRepo{Name: "testuser/testrepo"},
				Payload: map[string]interface{}{
					"commits": []interface{}{"commit1", "commit2"},
				},
			},
			expected: "Pushed 2 commits to testuser/testrepo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := pushEvent{event: tt.event}
			if got := e.Summarize(); got != tt.expected {
				t.Errorf("expected: %q, got: %q", tt.expected, got)
			}
		})
	}
}

func TestCreateEvent_Summarize(t *testing.T) {
	tests := []struct {
		name     string
		event    Event
		expected string
	}{
		{
			name: "valid create event",
			event: Event{
				Type: CreateEvent,
				Repo: EventRepo{Name: "testuser/testrepo"},
				Payload: map[string]interface{}{
					"ref_type": "branch",
				},
			},
			expected: "Created branch testuser/testrepo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := createEvent{event: tt.event}
			if got := e.Summarize(); got != tt.expected {
				t.Errorf("expected: %q, got: %q", tt.expected, got)
			}
		})
	}
}

func TestPullRequestEvent_Summarize(t *testing.T) {
	tests := []struct {
		name     string
		event    Event
		expected string
	}{
		{
			name: "valid pull request event",
			event: Event{
				Type: PullRequestEvent,
				Repo: EventRepo{Name: "testuser/testrepo"},
				Payload: map[string]interface{}{
					"action": "opened",
				},
			},
			expected: "Pull request opened testuser/testrepo",
		},
		{
			name: "valid pull request event with underscore separated action",
			event: Event{
				Type: PullRequestEvent,
				Repo: EventRepo{Name: "testuser/testrepo"},
				Payload: map[string]interface{}{
					"action": "review_request_removed",
				},
			},
			expected: "Pull request review request removed testuser/testrepo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := pullRequestEvent{event: tt.event}
			if got := e.Summarize(); got != tt.expected {
				t.Errorf("expected: %q, got: %q", tt.expected, got)
			}
		})
	}
}

func TestNotImplementedEvent_Summarize(t *testing.T) {
	tests := []struct {
		name     string
		event    Event
		expected string
	}{
		{
			name: "not implemented event",
			event: Event{
				Type: "UnknownEvent",
			},
			expected: fmt.Sprintf("%q is not implemented", "UnknownEvent"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := notImplementedEvent{event: tt.event}
			if got := e.Summarize(); got != tt.expected {
				t.Errorf("expected: %q, got: %q", tt.expected, got)
			}
		})
	}
}

func TestCommitComment_Summarize(t *testing.T) {
	tests := []struct {
		name     string
		event    Event
		expected string
	}{
		{
			name: "valid create event",
			event: Event{
				Type: CommitComment,
				Repo: EventRepo{Name: "testuser/testrepo"},
				Payload: map[string]interface{}{
					"action": "created",
				},
			},
			expected: "Commented on commit in testuser/testrepo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := commitCommentEvent{event: tt.event}
			if got := e.Summarize(); got != tt.expected {
				t.Errorf("expected: %q, got: %q", tt.expected, got)
			}
		})
	}
}

func TestIssueComment_Summarize(t *testing.T) {
	tests := []struct {
		name     string
		event    Event
		expected string
	}{
		{
			name: "valid create event",
			event: Event{
				Type: IssueCommentEvent,
				Repo: EventRepo{Name: "testuser/testrepo"},
				Payload: map[string]interface{}{
					"action": "created",
				},
			},
			expected: "Comment created on issue in testuser/testrepo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := issueCommentEvent{event: tt.event}
			if got := e.Summarize(); got != tt.expected {
				t.Errorf("expected: %q, got: %q", tt.expected, got)
			}
		})
	}
}

func TestDelete_Summarize(t *testing.T) {
	tests := []struct {
		name     string
		event    Event
		expected string
	}{
		{
			name: "valid create event",
			event: Event{
				Type: DeleteEvent,
				Repo: EventRepo{Name: "testuser/testrepo"},
				Payload: map[string]interface{}{
					"ref":      "feature/branch",
					"ref_type": "branch",
				},
			},
			expected: "Deleted branch feature/branch from testuser/testrepo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := deleteEvent{event: tt.event}
			if got := e.Summarize(); got != tt.expected {
				t.Errorf("expected: %q, got: %q", tt.expected, got)
			}
		})
	}
}

func TestFork_Summarize(t *testing.T) {
	tests := []struct {
		name     string
		event    Event
		expected string
	}{
		{
			name: "valid create event",
			event: Event{
				Type: ForkEvent,
				Repo: EventRepo{Name: "testuser/testrepo-fork"},
				Payload: map[string]interface{}{
					"forkee": map[string]interface{}{
						"full_name": "testuser/testrepo",
					},
				},
			},
			expected: "Forked testuser/testrepo to testuser/testrepo-fork",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := forkEvent{event: tt.event}
			if got := e.Summarize(); got != tt.expected {
				t.Errorf("expected: %q, got: %q", tt.expected, got)
			}
		})
	}
}
