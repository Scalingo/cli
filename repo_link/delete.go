package repo_link

import (
	"fmt"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func Delete(app string) error {
	if app == "" {
		return errgo.New("no app defined")
	}

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	// Get RepoLink of App
	repoLink, err := c.ScmRepoLinkShow(app)
	if err != nil {
		return errgo.Mask(err)
	}

	err = c.ScmRepoLinkDelete(app, repoLink.ID)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Printf("Current repo link has been deleted from app '%s'.\n", app)
	return nil
}
