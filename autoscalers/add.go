package autoscalers

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
	errgo "gopkg.in/errgo.v1"
)

func Add(app string, params scalingo.AutoscalerAddParams) error {
	c := config.ScalingoClient()
	autoscaler, err := c.AutoscalerAdd(app, params)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Autoscaler created on", app, "for", autoscaler.ContainerType, "containers")
	return nil
}
