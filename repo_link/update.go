package repo_link

import (
	"fmt"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo"
)

func Update(app string, repoLinkID string, params scalingo.ScmRepoLinkParams) error {
	if app == "" {
		return errgo.New("no app defined")
	}

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	check, err := checkRepoLinkExist(c, app, repoLinkID)
	if err != nil {
		return errgo.Mask(err)
	}
	if check == false {
		return errgo.Newf("RepoLink '%s' doesn't exist for app '%s'", repoLinkID, app)
	}

	_, err = c.ScmRepoLinkUpdate(app, repoLinkID, params)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Printf("RepoLink '%s' has been updated for app '%s'.\n", repoLinkID, app)
	return nil
}
