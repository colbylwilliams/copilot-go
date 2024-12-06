// Package copilot-go/github is a companion Go library for use with [copilot-go].
//
// It provides a set of utilities for working with the GitHub API.
//
// [copilot-go](https://pkg.go.dev/github.com/colbylwilliams/copilot-go)
package github

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/colbylwilliams/copilot-go"
	gh "github.com/google/go-github/v67/github"
	"golang.org/x/oauth2"
)

type ReferenceResource interface{}

func NewClient(ctx context.Context) (*gh.Client, error) {
	token := copilot.GetGetHubToken(ctx)
	if token == "" {
		return nil, fmt.Errorf("no github token found in context")
	}

	client := gh.NewClient(nil)
	if envURL := os.Getenv("GITHUB_API_URL"); envURL != "" {
		client.BaseURL, _ = url.Parse(envURL + "/")
	}
	return client.WithAuthToken(token), nil
}

func NewAppClient(ctx context.Context, cfg *copilot.Config) (*gh.Client, error) {
	if cfg.GitHubAppClientID == "" || string(cfg.GitHubAppPrivateKey) == "" {
		return nil, fmt.Errorf("no github app client id or private key in config")
	}
	ats, err := NewApplicationTokenSource(cfg.GitHubAppClientID, cfg.GitHubAppPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error creating application token source: %w", err)
	}

	httpClient := oauth2.NewClient(context.Background(), ats)

	client := gh.NewClient(httpClient)
	if envURL := os.Getenv("GITHUB_API_URL"); envURL != "" {
		client.BaseURL, _ = url.Parse(envURL + "/")
	}

	return client, nil
}
