package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/autoscalers"
	"github.com/Scalingo/cli/cmd/autocomplete"
	scalingo "github.com/Scalingo/go-scalingo"
	"github.com/urfave/cli"
)

var (
	autoscalersListCommand = cli.Command{
		Name:        "autoscalers",
		Category:    "Autoscalers",
		Usage:       "List the autoscalers of an application",
		Flags:       []cli.Flag{appFlag},
		Description: "List all the autoscalers of an application and display information about them.",
		Before:      AuthenticateHook,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "autoscalers")
				return
			}

			err := autoscalers.List(appdetect.CurrentApp(c))
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "autoscalers")
		},
	}

	autoscalersAddCommand = cli.Command{
		Name:     "autoscalers-add",
		Category: "Autoscalers",
		Usage:    "Add an autoscaler to an application",
		Flags: []cli.Flag{appFlag,
			cli.StringFlag{Name: "container-type, c", Usage: "Specify the container type affected by the autoscaler"},
			cli.StringFlag{Name: "metric, m", Usage: "Specify the metric you want the autoscaling to apply on"},
			cli.Float64Flag{Name: "target, t", Usage: "Target value for the metric the autoscaler will maintain"},
			cli.IntFlag{Name: "min-containers", Usage: "lower limit the autoscaler will never scale below"},
			cli.IntFlag{Name: "max-containers", Usage: "upper limit the autoscaler will never scale above"},
		},
		Description: `Add an autoscaler to an application. It will automatically scale the application up or down depending on the target defined for the given metric.

   All options are mandatory.

   Example
     scalingo --app my-app autoscaler-add --container-type web --metric cpu --target 0.75 --min-containers 1 --max-containers 3
		`,
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			if !isValidAutoscalerAddOpts(c) {
				err := cli.ShowCommandHelp(c, "autoscalers-add")
				if err != nil {
					errorQuit(err)
				}
				return
			}

			currentApp := appdetect.CurrentApp(c)
			err := autoscalers.Add(currentApp, scalingo.AutoscalerAddParams{
				ContainerType: c.String("c"),
				Metric:        c.String("m"),
				Target:        c.Float64("t"),
				MinContainers: c.Int("min-containers"),
				MaxContainers: c.Int("max-containers"),
			})
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			// TODO
		},
	}

	/*
		autoscalersRemoveCommand = cli.Command{
			Name:        "autoscalers-remove",
			Category:    "Autoscalers",
			Usage:       "Revoke permission to collaborate on an application",
			Flags:       []cli.Flag{appFlag},
			Description: "Revoke the invitation of collaboration to the given email\n    scalingo -a myapp autoscalers-remove user@example.com",
			Before:      AuthenticateHook,
			Action: func(c *cli.Context) {
				currentApp := appdetect.CurrentApp(c)
				if len(c.Args()) != 1 {
					cli.ShowCommandHelp(c, "autoscalers-remove")
				} else {
					err := autoscalers.Remove(currentApp, c.Args()[0])
					if err != nil {
						errorQuit(err)
					}
				}
			},
			BashComplete: func(c *cli.Context) {
				autocomplete.CmdFlagsAutoComplete(c, "autoscalers-remove")
				autocomplete.AutoscalersRemoveAutoComplete(c)
			},
		}
	*/
)

func isValidAutoscalerAddOpts(c *cli.Context) bool {
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
