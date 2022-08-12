package update

import (
	"context"
	"fmt"
	"strings"

	"gopkg.in/errgo.v1"

	githubClient "github.com/google/go-github/v45/github"

	"github.com/Scalingo/cli/services/github"
)

func printRelease(release githubClient.RepositoryRelease) {
	fmt.Printf("Changelog of the version %v\n\n", release.GetTagName())
	fmt.Printf("%v\n\n", strings.ReplaceAll(*release.Body, "\\r\\n", "\r\n"))
}

func ShowLastChangelog() error {
	ctx := context.Background()

	client := github.New()

	cliLastRelease, err := client.GetLatestRelease(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get last CLI release")
	}
	printRelease(*cliLastRelease)

	return nil
}

func ShowChangelog(cacheVersion, newVersion string) error {
	cliLastReleases, err := getLastReleases()
	if err != nil {
		return errgo.Notef(err, "fail to get last CLI release")
	}

	beginToSowRelease := false
	// Read it from the oldest to the most recent.
	for i := len(cliLastReleases) - 1; i >= 0; i-- {
		release := cliLastReleases[i]
		if beginToSowRelease {
			printRelease(*release)
		}

		// Detect the release from the cache, we can show changelog from next release.
		if cacheVersion == release.GetTagName() {
			beginToSowRelease = true
		}
	}

	return nil
}

// getLastReleases get last 10 releases.
// It returns releases in map, the key is the tag, the value is the body of the release.
// The releases are ordered by the last in first.
func getLastReleases() ([]*githubClient.RepositoryRelease, error) {
	ctx := context.Background()

	client := github.New()

	// Only show the last 10 releases in maximum
	cliReleases, err := client.ListReleases(ctx, &githubClient.ListOptions{PerPage: 10})
	if err != nil {
		return nil, errgo.Notef(err, "fail to list CLI releases")
	}

	return cliReleases, nil
}
