package github

import (
	"context"

	"github.com/google/go-github/v45/github"

	"github.com/Scalingo/go-scalingo/v4/debug"
)

const (
	owner = "scalingo"
	repo  = "cli"
)

type client struct {
	githubRepoService *github.RepositoriesService
}

func (c client) ListReleases(ctx context.Context, opts *github.ListOptions) ([]*github.RepositoryRelease, error) {
	repoReleases, githubResponse, error := c.githubRepoService.ListReleases(ctx, owner, repo, opts)
	debug.Printf("github response: %#v\n", githubResponse)
	return repoReleases, error
}

func (c client) GetLatestRelease(ctx context.Context) (*github.RepositoryRelease, error) {
	repoRelease, githubResponse, error := c.githubRepoService.GetLatestRelease(ctx, owner, repo)
	debug.Printf("github response: %#v\n", githubResponse)
	return repoRelease, error
}

func New() client {
	githubClient := github.NewClient(nil)
	return client{
		githubRepoService: githubClient.Repositories,
	}
}
