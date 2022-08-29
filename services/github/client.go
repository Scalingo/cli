package github

import (
	"context"
	"net/http"
	"time"

	"github.com/google/go-github/v47/github"

	"github.com/Scalingo/go-scalingo/v4/debug"
)

const (
	owner = "scalingo"
	repo  = "cli"
)

type Client interface {
	GetLatestRelease(ctx context.Context) (*github.RepositoryRelease, error)
}

type client struct {
	githubRepositoriesService *github.RepositoriesService
}

func NewClient() Client {
	return client{
		githubRepositoriesService: github.NewClient(&http.Client{
			Timeout: 5 * time.Second,
		}).Repositories,
	}
}

func (c client) GetLatestRelease(ctx context.Context) (*github.RepositoryRelease, error) {
	repoRelease, githubResponse, err := c.githubRepositoriesService.GetLatestRelease(ctx, owner, repo)
	debug.Printf("GitHub response: %#v\n", githubResponse)
	return repoRelease, err
}
