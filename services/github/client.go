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

type Client struct {
	githubRepoService *github.RepositoriesService
}

func (c Client) ListReleases(ctx context.Context, opts *github.ListOptions) ([]*github.RepositoryRelease, error) {
	repoReleases, githubResponse, err := c.githubRepoService.ListReleases(ctx, owner, repo, opts)
	debug.Printf("GitHub response: %#v\n", githubResponse)
	return repoReleases, err
}

func (c Client) GetLatestRelease(ctx context.Context) (*github.RepositoryRelease, error) {
	repoRelease, githubResponse, err := c.githubRepoService.GetLatestRelease(ctx, owner, repo)
	debug.Printf("github response: %#v\n", githubResponse)
	return repoRelease, err
}

func NewClient() Client {
	githubClient := github.NewClient(nil)
	return Client{
		githubRepoService: githubClient.Repositories,
	}
}
