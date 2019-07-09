package alerts

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
)

func Add(app string, params scalingo.AlertAddParams) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	a, err := c.AlertAdd(app, params)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Status("Alert created for the container type", a.ContainerType)
	return nil
}
