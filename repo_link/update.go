package repo_link

import (
	"fmt"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo"
)

func Update(app string, params scalingo.ScmRepoLinkParams) error {
	if app == "" {
		return errgo.New("no app defined")
	}

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	_, err = c.ScmRepoLinkUpdate(app, params)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Printf("RepoLink has been updated for app '%s'.\n", app)
	return nil
}
