package alerts

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
)

func Add(app string, params scalingo.AlertParams) error {
	c := config.ScalingoClient()
	a, err := c.AlertAdd(app, params)

	if err != nil {
		return errgo.Mask(err)
	}

	io.Status("Alert created for the container type", a.ContainerType)
	io.Info("http://")
	return nil
}
