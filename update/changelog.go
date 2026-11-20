package update

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	gogithub "github.com/google/go-github/v88/github"

	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/services/github"
	"github.com/Scalingo/go-utils/errors/v3"
)

const compareURL = "https://github.com/Scalingo/cli/compare/%s...%s"

// fullChangelogLine matches the compare link GitHub adds at the end of each release body
var fullChangelogLine = regexp.MustCompile(`(?m)^\*\*Full Changelog\*\*:.*$`)

// ShowLastChangelog shows the changelog of the last release
func ShowLastChangelog(ctx context.Context) error {
	return ShowChangelogSince(ctx, "")
}

// ShowChangelogSince shows the changelog of every release published after the version since.
// With an empty or unknown version, it shows the changelog of the last release.
func ShowChangelogSince(ctx context.Context, since string) error {
	githubClient, err := github.NewClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "initialize GitHub client")
	}

	var releases []*gogithub.RepositoryRelease
	if since == "" {
		cliLastRelease, err := githubClient.GetLatestRelease(ctx)
		if err != nil {
			return errors.Wrapf(ctx, err, "get last CLI release")
		}
		releases = []*gogithub.RepositoryRelease{cliLastRelease}
	} else {
		releases, err = githubClient.ListReleases(ctx)
		if err != nil {
			return errors.Wrapf(ctx, err, "list CLI releases")
		}
	}

	changelog := changelogSince(releases, since)
	if changelog == "" {
		io.Warning("GitHub last release is empty")
		return nil
	}

	fmt.Print(changelog)

	return nil
}

// changelogSince renders the changelog of the releases published after the version since.
// releases must be sorted newest first, as GitHub lists them. Drafts and pre-releases are skipped.
func changelogSince(releases []*gogithub.RepositoryRelease, since string) string {
	published := make([]*gogithub.RepositoryRelease, 0, len(releases))
	for _, release := range releases {
		if release.GetDraft() || release.GetPrerelease() {
			continue
		}
		published = append(published, release)
	}
	if len(published) == 0 {
		return ""
	}

	newer := published[:1]
	for i, release := range published {
		if release.GetTagName() == since {
			if i > 0 {
				newer = published[:i]
			}
			break
		}
	}

	if len(newer) == 1 {
		release := newer[0]
		if release.GetBody() == "" {
			return ""
		}
		return fmt.Sprintf("Changelog of the version %v\n\n%v\n\n", release.GetTagName(), releaseBody(release))
	}

	var changelog strings.Builder
	for _, release := range newer {
		body := strings.TrimSpace(fullChangelogLine.ReplaceAllString(releaseBody(release), ""))
		fmt.Fprintf(&changelog, "Changelog of the version %v\n\n%v\n\n", release.GetTagName(), body)
	}
	fmt.Fprintf(&changelog, "**Full Changelog**: "+compareURL+"\n\n", since, newer[0].GetTagName())

	return changelog.String()
}

func releaseBody(release *gogithub.RepositoryRelease) string {
	return strings.ReplaceAll(release.GetBody(), "\\r\\n", "\r\n")
}
