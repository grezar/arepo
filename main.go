package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v26/github"
	"golang.org/x/oauth2"
)

const (
	version = "0.1.1"
)

var (
	displayVersion bool
	isPrivate      bool
	accessToken    string
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Invalid number of arguments:", len(os.Args))
		os.Exit(1)
	}
	os.Exit(run())
}

func run() int {
	flag.BoolVar(&displayVersion, "version", false, "Display version")
	flag.BoolVar(&isPrivate, "private", false, "Create a new repository with private")
	flag.StringVar(&accessToken, "token", "", "GitHub access token")
	flag.Parse()
	if displayVersion {
		fmt.Println(version)
		return 0
	}
	if err := doCommand(); err != nil {
		fmt.Println("Failed to create a new repository: ", err)
		return 1
	}
	return 0
}

func doCommand() error {
	client, err := newGitHubClient()
	if err != nil {
		return err
	}
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

func newGitHubClient() (*github.Client, error) {
	ctx := context.Background()
	token, err := getAccessToken()
	if err != nil {
		return nil, err
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), nil
}

func getAccessToken() (string, error) {
	if accessToken != "" {
		return accessToken, nil
	}
	if t := os.Getenv("GITHUB_ACCESS_TOKEN"); t != "" {
		return t, nil
	}
	areporcPath := fmt.Sprintf("%s/.areporc", os.Getenv("HOME"))
	if _, err := os.Stat(areporcPath); !os.IsNotExist(err) {
		file, err := os.Open(areporcPath)
		if err != nil {
			return "", err
		}
		defer file.Close()
		buf := make([]byte, 64)
		for {
			n, err := file.Read(buf)
			if n == 0 {
				break
			}
			if err != nil {
				return "", err
			}
			return strings.TrimSuffix(string(buf[:n]), "\n"), nil
		}
	}
	return "", nil
}
