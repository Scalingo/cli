package github

import (
	"context"

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
	githubRepoService *github.RepositoriesService
}

func NewClient() Client {
	githubClient := github.NewClient(nil)
	return client{
		githubRepoService: githubClient.Repositories,
	}
}

func (c client) GetLatestRelease(ctx context.Context) (*github.RepositoryRelease, error) {
	repoRelease, githubResponse, err := c.githubRepoService.GetLatestRelease(ctx, owner, repo)
	debug.Printf("GitHub response: %#v\n", githubResponse)
	return repoRelease, err
}
