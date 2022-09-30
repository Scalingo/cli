package git

import (
	"context"

	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func Show(ctx context.Context, appName string) error {
	url, err := getGitEndpoint(ctx, appName)
	if err != nil {
		return errgo.Notef(err, "fail to get the Git endpoint of this app")
	}

	io.Status("Your application", appName, "Git remote is:", url)
	return nil
}

func getGitEndpoint(ctx context.Context, appName string) (string, error) {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return "", errgo.Notef(err, "fail to get scalingo API client")
	}
	app, err := c.AppsShow(ctx, appName)
	if err != nil {
		return "", errgo.Notef(err, "fail to get application information")
	}

	return app.GitURL, nil
}
