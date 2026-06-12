package github

import (
	"context"
	"net/http"
	"time"

	"github.com/google/go-github/v88/github"

	"github.com/Scalingo/go-scalingo/v11/debug"
	"github.com/Scalingo/go-utils/errors/v3"
)

type Client interface {
	GetLatestRelease(ctx context.Context) (*github.RepositoryRelease, error)
}

type client struct {
	githubRepositoriesService *github.RepositoriesService
}

func NewClient(ctx context.Context) (Client, error) {
	githubClient, err := github.NewClient(github.WithHTTPClient(&http.Client{
		Timeout: 5 * time.Second,
	}))
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "initialize GitHub HTTP client")
	}

	return client{
		githubRepositoriesService: githubClient.Repositories,
	}, nil
}

func (c client) GetLatestRelease(ctx context.Context) (*github.RepositoryRelease, error) {
	latestRelease, githubResponse, err := c.githubRepositoriesService.GetLatestRelease(ctx, "Scalingo", "cli")
	if githubResponse != nil && githubResponse.Body != nil {
		defer githubResponse.Body.Close()
	}

	debug.Printf("GitHub response: %#v\n", githubResponse)

	if err != nil {
		return nil, errors.Wrapf(ctx, err, "fail to get the latest release of the Scalingo/cli repository")
	}

	return latestRelease, nil
}
