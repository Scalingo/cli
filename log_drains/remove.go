package log_drains

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func Remove(app string, URL string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	err = c.LogDrainRemove(app, URL)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Status("The log drain:", URL, "has been deleted")
	return nil
}
