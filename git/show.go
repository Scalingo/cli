package git

import (
	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func Show(appName string) error {
	url, err := getGitEndpoint(appName)
	if err != nil {
		return errgo.Notef(err, "fail to get the Git endpoint of this app")
	}

	io.Status("Your application", appName, "Git remote is:", url)
	return nil
}

func getGitEndpoint(appName string) (string, error) {
	c, err := config.ScalingoClient()
	if err != nil {
		return "", errgo.Notef(err, "fail to get scalingo API client")
	}
	app, err := c.AppsShow(appName)
	if err != nil {
		return "", errgo.Notef(err, "fail to get application information")
	}

	return app.GitUrl, nil
}
