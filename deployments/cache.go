package deployments

import (
	"github.com/Scalingo/cli/config"
	"gopkg.in/errgo.v1"
)

func ResetCache(app string) error {
	c := config.ScalingoClient()
	err := c.DeploymentCacheReset(app)

	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}
