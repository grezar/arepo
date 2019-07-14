package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v26/github"
	"golang.org/x/oauth2"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Invalid number of arguments (expected 2, got %d)\n", len(os.Args))
		os.Exit(1)
	}
	os.Exit(run())
}

func run() int {
	if err := doCommand(); err != nil {
		fmt.Println("Failed to create a new repository: ", err)
		return 1
	}
	return 0
}

func doCommand() error {
	client := newGitHubClient()
	ctx := context.Background()
	repo := &github.Repository{
		Name:    github.String(os.Args[1]),
		Private: github.Bool(false),
	}
	r, _, err := client.Repositories.Create(ctx, "", repo)
	if err != nil {
		return err
	}
	fmt.Println(*r.CloneURL)
	return nil
}

func newGitHubClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: getAccessToken()},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func getAccessToken() string {
	return os.Getenv("GITHUB_ACCESS_TOKEN")
}
