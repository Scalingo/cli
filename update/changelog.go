package update

import (
	"context"
	"fmt"
	"strings"

	githubClient "github.com/google/go-github/v47/github"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/services/github"
)

func printRelease(release githubClient.RepositoryRelease) {
	fmt.Printf("Changelog of the version %v\n\n", release.GetTagName())
	fmt.Printf("%v\n\n", strings.ReplaceAll(*release.Body, "\\r\\n", "\r\n"))
}

func ShowLastChangelog() error {
	ctx := context.Background()

	client := github.NewClient()

	cliLastRelease, err := client.GetLatestRelease(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get last CLI release")
	}
	printRelease(*cliLastRelease)

	return nil
}
