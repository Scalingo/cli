package repolink

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
)

func Update(app string, params scalingo.SCMRepoLinkParams) error {
	if app == "" {
		return errgo.New("no app defined")
	}

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	_, err = c.SCMRepoLinkUpdate(app, params)
	if err != nil {
		return errgo.Notef(err, "fail to update repo link")
	}

	io.Statusf("RepoLink has been updated for app '%s'.\n", app)
	return nil
}
