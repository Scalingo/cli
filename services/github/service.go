package github

import (
	"context"

	"github.com/google/go-github/v45/github"
)

type Service interface {
	ListReleases(ctx context.Context, opts *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error)
	GetLatestRelease(ctx context.Context) (*github.RepositoryRelease, *github.Response, error)
}
