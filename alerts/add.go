package alerts

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
)

func Add(app string, params scalingo.AlertAddParams) error {
	a, err := config.ScalingoClient().AlertAdd(app, params)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Status("Alert created for the container type", a.ContainerType)
	io.Info("http://")
	return nil
}
