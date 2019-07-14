package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/go-github/v26/github"
	"golang.org/x/oauth2"
)

var (
	isPrivate bool
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Invalid number of arguments:", len(os.Args))
		os.Exit(1)
	}
	os.Exit(run())
}

func run() int {
	flag.BoolVar(&isPrivate, "private", false, "Create a new repository with private")
	flag.Parse()
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
		Name:    github.String(flag.Arg(0)),
		Private: github.Bool(isPrivate),
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
