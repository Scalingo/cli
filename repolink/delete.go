package repolink

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func Delete(app string) error {
	if app == "" {
		return errgo.New("no app defined")
	}

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	err = c.SCMRepoLinkDelete(app)
	if err != nil {
		return errgo.Notef(err, "fail to delete repo link")
	}

	io.Statusf("Current repo link has been deleted from app '%s'.\n", app)
	return nil
}
