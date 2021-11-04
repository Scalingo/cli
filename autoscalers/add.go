package autoscalers

import (
	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v4"
)

func Add(app string, params scalingo.AutoscalerAddParams) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	autoscaler, err := c.AutoscalerAdd(app, params)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Autoscaler created on", app, "for", autoscaler.ContainerType, "containers")
	return nil
}
