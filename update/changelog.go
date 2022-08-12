package update

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v45/github"
	"gopkg.in/errgo.v1"
)

func ShowLastChangelog() error {
	ctx := context.Background()

	client := github.NewClient(nil)
	repoService := client.Repositories

	cliLastRelease, _, err := repoService.GetLatestRelease(ctx, "scalingo", "cli")
	if err != nil {
		return errgo.Notef(err, "fail to get last CLI release")
	}
	fmt.Printf("%v\n\n", strings.ReplaceAll(*cliLastRelease.Body, "\\r\\n", "\r\n"))

	return nil
}
