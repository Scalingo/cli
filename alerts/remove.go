package alerts

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func Remove(app, id string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	err = c.AlertRemove(app, id)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Status("The alert has been deleted")
	return nil
}
