package log_drains

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
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

func findLogDrainsBeforeRemove(c *scalingo.Client, app string, URL string) (scalingo.LogDrain, error) {
	drains, err := c.LogDrainsList(app)
	if err != nil {
		return scalingo.LogDrain{}, errgo.Mask(err)
	}

	for _, d := range drains {
		if d.URL == URL {
			return d, nil
		}
	}
	return scalingo.LogDrain{}, errgo.New("There is no such log drain, please ensure you've added it correctly.")
}
