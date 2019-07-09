package deployments

import (
	"github.com/Scalingo/cli/config"
	"gopkg.in/errgo.v1"
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
