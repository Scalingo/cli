package log_drains

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
)

func Add(app string, params scalingo.LogDrainAddParams) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client to add a log drain")
	}

	d, err := c.LogDrainAdd(app, params)
	if err != nil {
		return errgo.Notef(err, "fail to add drain to the application")
	}

	io.Status("Log drain", d.Drain.URL, "has been added to the application")
	return nil
}
