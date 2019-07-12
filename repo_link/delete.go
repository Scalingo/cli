package repo_link

import (
	"fmt"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo"
)

func Delete(app string, repoLinkID string) error {
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

	err = c.ScmRepoLinkDelete(app, repoLinkID)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Printf("RepoLink '%s' has been deleted from app '%s'.\n", repoLinkID, app)
	return nil
}

func checkRepoLinkExist(c *scalingo.Client, app, repoLinkID string) (bool, error) {
	repoLink, err := c.ScmRepoLinkShow(app)
	if err != nil {
		return false, errgo.Mask(err, errgo.Any)
	}

	if repoLink != nil && repoLink.ID == repoLinkID {
		return true, nil
	}

	return false, nil
}
