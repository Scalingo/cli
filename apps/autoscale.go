package apps

import (
	"fmt"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo"
	errgo "gopkg.in/errgo.v1"
)

type AutoscaleRes struct {
	Containers []scalingo.ContainerType `json:"containers"`
}

type AutoscaleOpts struct {
	App           string
	ContainerType string
	Metric        string
	Target        float64
	MinContainers int
	MaxContainers int
}

func Autoscale(opts AutoscaleOpts) error {
	autoscalerParams := &scalingo.AppsAutoscalerParams{
		Autoscaler: scalingo.Autoscaler{
			ContainerType: opts.ContainerType,
			Metric:        opts.Metric,
			Target:        opts.Target,
			MinContainers: opts.MinContainers,
			MaxContainers: opts.MaxContainers,
		},
	}

	c := config.ScalingoClient()
	res, err := c.AppsAutoscaler(app, autoscalerParams)
	if err != nil {
		return errgo.Mask(err)
	}
	defer res.Body.Close()

	fmt.Println("Autoscaler created for your application.")
	return nil
}
