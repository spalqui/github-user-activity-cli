package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spalqui/github-user-activity-cli/services/github"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a GitHub username")
	}

	username := os.Args[1]

	// Create a new HTTP client for GitHub API requests
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create a new GitHub service instance
	gitHubService := github.NewService(httpClient)

	// Get user events for a specific GitHub username
	events, err := gitHubService.GetUserEvents(username)
	if err != nil {
		log.Fatalf("Unable to get user activity for %q: %v", username, err)
	}

	fmt.Printf("%+v\n", events)
}
