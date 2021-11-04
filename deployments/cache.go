package deployments

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func ResetCache(app string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	err = c.DeploymentCacheReset(app)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}
