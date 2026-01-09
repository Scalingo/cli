package github

import (
	"context"
	"net/http"
	"time"

	"github.com/google/go-github/v47/github"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v9/debug"
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
	latestRelease, githubResponse, err := c.githubRepositoriesService.GetLatestRelease(ctx, "Scalingo", "cli")
	if githubResponse != nil && githubResponse.Body != nil {
		defer githubResponse.Body.Close()
	}

	debug.Printf("GitHub response: %#v\n", githubResponse)

	if err != nil {
		return nil, errgo.Notef(err, "fail to get the latest release of the Scalingo/cli repository")
	}

	return latestRelease, nil
}
