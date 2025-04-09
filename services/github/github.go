package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "https://api.github.com"

type Service struct {
	baseURL    string
	httpClient *http.Client
}

func NewService(httpClient *http.Client) *Service {
	return &Service{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (s *Service) GetUserEvents(username string) ([]Event, error) {
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
