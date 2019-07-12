package repo_link

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo"
)

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
