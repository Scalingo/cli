package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/urfave/cli"
)

var (
	ScaleCommand = cli.Command{
		Name:      "scale",
		ShortName: "s",
		Category:  "App Management",
		Flags: []cli.Flag{appFlag,
			cli.BoolFlag{Name: "synchronous, s", Usage: "Do the scaling synchronously", EnvVar: ""},
			cli.StringFlag{Name: "container-type, c", Usage: "Specify the container type affected by the autoscaler"},
			cli.StringFlag{Name: "metric, m", Usage: "Specify the metric you want the autoscaling to apply on"},
			cli.Float64Flag{Name: "target, t", Usage: "Target value for the metric the autoscaler will maintain"},
			cli.IntFlag{Name: "min-containers", Usage: "lower limit the autoscaler will never scale below"},
			cli.IntFlag{Name: "max-containers", Usage: "upper limit the autoscaler will never scale above"},
		},
		Usage: "Scale your application instantly",
		Description: `Scale your application processes.
   Example
     'scalingo --app my-app scale web:2 worker:1'
     'scalingo --app my-app scale web:1 worker:0'
     'scalingo --app my-app scale web:1:XL'
     'scalingo --app my-app scale web:+1 worker:-1'
     'scalingo --app my-app scale --container-type web --metric cpu --target 0.75 --min-containers 1 --max-containers 3'
     `,
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			if !isValidAutoscalerOpts(c) && !isValidScaleOpts(c) {
				err := cli.ShowCommandHelp(c, "scale")
				if err != nil {
					errorQuit(err)
				}
				return
			}

			currentApp := appdetect.CurrentApp(c)

			if len(c.Args()) > 0 {
				err := apps.Scale(currentApp, c.Bool("s"), c.Args())
				if err != nil {
					errorQuit(err)
				}
				return
			}

			opts := apps.AutoscaleOpts{
				App:           currentApp,
				ContainerType: c.String("container-type"),
				Metric:        c.String("metric"),
				Target:        c.Float64("target"),
				MinContainers: c.Int("min-containers"),
				MaxContainers: c.Int("max-containers"),
			}

			err := apps.Autoscale(opts)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "scale")
			autocomplete.ScaleAutoComplete(c)
		},
	}
)

func isValidAutoscalerOpts(c *cli.Context) bool {
	if len(c.Args()) > 0 {
		return false
	}
	for _, opt := range []string{
		"container-type", "metric", "target", "min-containers", "max-containers",
	} {
		if !c.IsSet(opt) {
			return false
		}
	}
	return true
}

func isValidScaleOpts(c *cli.Context) bool {
	if len(c.Args()) == 0 {
		return false
	}
	return true
}
