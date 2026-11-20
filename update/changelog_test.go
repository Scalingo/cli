package update

import (
	"testing"

	"github.com/google/go-github/v88/github"
	"github.com/stretchr/testify/assert"
)

func release(tag, body string) *github.RepositoryRelease {
	return &github.RepositoryRelease{TagName: github.Ptr(tag), Body: github.Ptr(body)}
}

func TestChangelogSince(t *testing.T) {
	// GitHub lists the releases newest first
	releases := []*github.RepositoryRelease{
		release("1.43.3", "* fix: third\r\n\r\n**Full Changelog**: https://github.com/Scalingo/cli/compare/1.43.2...1.43.3"),
		release("1.43.2", "* feat: second\r\n\r\n**Full Changelog**: https://github.com/Scalingo/cli/compare/1.43.1...1.43.2"),
		release("1.43.1", "* fix: first\r\n\r\n**Full Changelog**: https://github.com/Scalingo/cli/compare/1.43.0...1.43.1"),
		release("1.43.0", "* feat: zero"),
	}

	tests := map[string]struct {
		releases []*github.RepositoryRelease
		since    string
		want     string
	}{
		"without a version, it shows the last release as before": {
			releases: releases,
			since:    "",
			want: "Changelog of the version 1.43.3\n\n" +
				"* fix: third\r\n\r\n**Full Changelog**: https://github.com/Scalingo/cli/compare/1.43.2...1.43.3\n\n",
		},
		"it shows every release after the installed version and one compare link": {
			releases: releases,
			since:    "1.43.1",
			want: "Changelog of the version 1.43.3\n\n* fix: third\n\n" +
				"Changelog of the version 1.43.2\n\n* feat: second\n\n" +
				"**Full Changelog**: https://github.com/Scalingo/cli/compare/1.43.1...1.43.3\n\n",
		},
		"the installed version is the previous one": {
			releases: releases,
			since:    "1.43.2",
			want: "Changelog of the version 1.43.3\n\n" +
				"* fix: third\r\n\r\n**Full Changelog**: https://github.com/Scalingo/cli/compare/1.43.2...1.43.3\n\n",
		},
		"the installed version is not in the list, it shows the last release": {
			releases: releases,
			since:    "1.20.0",
			want: "Changelog of the version 1.43.3\n\n" +
				"* fix: third\r\n\r\n**Full Changelog**: https://github.com/Scalingo/cli/compare/1.43.2...1.43.3\n\n",
		},
		"drafts and pre-releases are skipped": {
			releases: append([]*github.RepositoryRelease{
				{TagName: github.Ptr("1.44.0-beta"), Body: github.Ptr("* beta"), Prerelease: github.Ptr(true)},
				{TagName: github.Ptr("1.44.0"), Body: github.Ptr("* draft"), Draft: github.Ptr(true)},
			}, releases...),
			since: "1.43.2",
			want: "Changelog of the version 1.43.3\n\n" +
				"* fix: third\r\n\r\n**Full Changelog**: https://github.com/Scalingo/cli/compare/1.43.2...1.43.3\n\n",
		},
		"no release": {
			releases: nil,
			since:    "1.43.1",
			want:     "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.want, changelogSince(test.releases, test.since))
		})
	}
}
