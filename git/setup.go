package git

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/debug"
	errgo "gopkg.in/errgo.v1"
	git "gopkg.in/src-d/go-git.v4"
	gitconfig "gopkg.in/src-d/go-git.v4/config"
)

type SetupParams struct {
	RemoteName string
}

func Setup(appName string, params SetupParams) error {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return errgo.Notef(err, "fail to initialize the Git repository")
	}

	url, err := getGitEndpoint(appName)
	if err != nil {
		return errgo.Notef(err, "fail to get the Git endpoint of this app")
	}
	debug.Println("Adding Git remote", url)

	_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
		Name: params.RemoteName,
		URLs: []string{url},
	})
	if err != nil {
		return errgo.Notef(err, "fail to add the Git remote")
	}

	io.Status("Successfully added the Git remote", params.RemoteName, "on", appName)
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
