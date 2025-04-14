package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	BaseURL = "https://api.github.com"

	CreateEvent      = "CreateEvent"
	PushEvent        = "PushEvent"
	PullRequestEvent = "PullRequestEvent"
)

type Service struct {
	baseURL    string
	httpClient *http.Client
}

type Option func(*Service)

func NewService(opts ...Option) *Service {
	s := &Service{
		baseURL: BaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func WithBaseURL(baseURL string) Option {
	return func(s *Service) {
		s.baseURL = baseURL
	}
}

func (s *Service) GetUserEventsSummary(username string) (string, error) {
	events, err := s.getUserEvents(username)
	if err != nil {
		return "", fmt.Errorf("error getting user events: %v", err)
	}
	return generateEventsSummary(events), nil
}

type Event struct {
	Type      string                 `json:"type"`
	Repo      EventRepo              `json:"repo"`
	Payload   map[string]interface{} `json:"payload"`
	CreatedAt string                 `json:"created_at"`
}

type EventRepo struct {
	Name string `json:"name"`
}

func (s *Service) getUserEvents(username string) ([]Event, error) {
	url := fmt.Sprintf("%s/users/%s/events", s.baseURL, username)
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error non ok status code: %s", resp.Status)
	}

	var events []Event
	err = json.NewDecoder(resp.Body).Decode(&events)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return events, nil
}

type eventTypeSummarizer interface {
	Summarize() string
}

func generateEventsSummary(events []Event) string {
	summary := strings.Builder{}
	summary.WriteString("Output:\n")
	for _, event := range events {
		var e eventTypeSummarizer
		switch event.Type {
		case CreateEvent:
			e = createEvent{event: event}
		case PushEvent:
			e = pushEvent{event: event}
		case PullRequestEvent:
			e = pullRequestEvent{event: event}
		default:
			e = notImplementedEvent{event: event}
		}
		summary.WriteString(fmt.Sprintf("- %s", e.Summarize()))
	}
	return summary.String()
}

type notImplementedEvent struct {
	event Event
}

func (e notImplementedEvent) Summarize() string {
	return fmt.Sprintf("%q is not implemented\n", e.event.Type)
}

type pushEvent struct {
	event Event
}

func (e pushEvent) Summarize() string {
	return fmt.Sprintf("Pushed %d commits to %s\n", len(e.event.Payload["commits"].([]interface{})), e.event.Repo.Name)
}

type createEvent struct {
	event Event
}

func (e createEvent) Summarize() string {
	return fmt.Sprintf("Created %s %s\n", e.event.Payload["ref_type"], e.event.Repo.Name)
}

type pullRequestEvent struct {
	event Event
}

func (e pullRequestEvent) Summarize() string {
	action := e.event.Payload["action"]
	action = strings.Join(strings.Split(action.(string), "_"), " ")
	return fmt.Sprintf("Pull request %s %s\n", e.event.Payload["action"], e.event.Repo.Name)
}
