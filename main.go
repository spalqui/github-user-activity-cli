package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spalqui/github-user-activity-cli/services/github"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a GitHub username")
	}

	username := os.Args[1]

	// Create a new GitHub service instance
	gitHubService := github.NewService()

	// Get user events for a specific GitHub username
	summary, err := gitHubService.GetUserEventsSummary(username)
	if err != nil {
		log.Fatalf("Unable to get user activity for %q: %v", username, err)
	}

	fmt.Print(summary)
}
