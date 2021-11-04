package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/autoscalers"
	"github.com/Scalingo/cli/cmd/autocomplete"
	scalingo "github.com/Scalingo/go-scalingo/v4"
)

var (
	autoscalersListCommand = cli.Command{
		Name:        "autoscalers",
		Category:    "Autoscalers",
		Usage:       "List the autoscalers of an application",
		Flags:       []cli.Flag{appFlag},
		Description: "List all the autoscalers of an application and display information about them.",
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
     scalingo --app my-app autoscalers-add --container-type web --metric cpu --target 0.75 --min-containers 1 --max-containers 3
		`,
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
			autocomplete.CmdFlagsAutoComplete(c, "autoscalers-add")
		},
	}

	autoscalersUpdateCommand = cli.Command{
		Name:     "autoscalers-update",
		Category: "Autoscalers",
		Usage:    "Update an autoscaler",
		Flags: []cli.Flag{appFlag,
			cli.StringFlag{Name: "container-type, c", Usage: "Specify the container type affected by the autoscaler"},
			cli.StringFlag{Name: "metric, m", Usage: "Specify the metric you want the autoscaling to apply on"},
			cli.Float64Flag{Name: "target, t", Usage: "Target value for the metric the autoscaler will maintain"},
			cli.IntFlag{Name: "min-containers", Usage: "lower limit the autoscaler will never scale below"},
			cli.IntFlag{Name: "max-containers", Usage: "upper limit the autoscaler will never scale above"},
			cli.BoolFlag{Name: "disabled, d", Usage: "disable/enable the given autoscaler"},
		},
		Description: `Update an autoscaler.

   The "container-type" option is mandatory.

   Example
     scalingo --app my-app autoscalers-update --container-type web --max-containers 5
     scalingo --app my-app autoscalers-update --container-type web --metric p95_response_time --target 67
		`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 0 || !c.IsSet("c") {
				err := cli.ShowCommandHelp(c, "autoscalers-update")
				if err != nil {
					errorQuit(err)
				}
				return
			}

			currentApp := appdetect.CurrentApp(c)
			params := scalingo.AutoscalerUpdateParams{}
			if c.IsSet("m") {
				m := c.String("m")
				params.Metric = &m
			}
			if c.IsSet("t") {
				t := c.Float64("t")
				params.Target = &t
			}
			if c.IsSet("min-containers") {
				min := c.Int("min-containers")
				params.MinContainers = &min
			}
			if c.IsSet("max-containers") {
				max := c.Int("max-containers")
				params.MaxContainers = &max
			}
			if c.IsSet("d") {
				d := c.Bool("d")
				params.Disabled = &d
			}
			err := autoscalers.Update(currentApp, c.String("c"), params)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "autoscalers-update")
		},
	}

	autoscalersEnableCommand = cli.Command{
		Name:     "autoscalers-enable",
		Category: "Autoscalers",
		Usage:    "Enable an autoscaler",
		Flags:    []cli.Flag{appFlag},
		Description: `Enable an autoscaler.

   Example
     scalingo --app my-app autoscalers-enable web
		`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				err := cli.ShowCommandHelp(c, "autoscalers-enable")
				if err != nil {
					errorQuit(err)
				}
				return
			}

			currentApp := appdetect.CurrentApp(c)
			disabled := false
			err := autoscalers.Update(currentApp, c.Args()[0], scalingo.AutoscalerUpdateParams{
				Disabled: &disabled,
			})
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "autoscalers-enable")
		},
	}

	autoscalersDisableCommand = cli.Command{
		Name:     "autoscalers-disable",
		Category: "Autoscalers",
		Usage:    "Disable an autoscaler",
		Flags:    []cli.Flag{appFlag},
		Description: `Disable an autoscaler.

   Example
     scalingo --app my-app autoscalers-disable web
		`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				err := cli.ShowCommandHelp(c, "autoscalers-disable")
				if err != nil {
					errorQuit(err)
				}
				return
			}

			currentApp := appdetect.CurrentApp(c)
			disabled := true
			err := autoscalers.Update(currentApp, c.Args()[0], scalingo.AutoscalerUpdateParams{
				Disabled: &disabled,
			})
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "autoscalers-disable")
		},
	}

	autoscalersRemoveCommand = cli.Command{
		Name:     "autoscalers-remove",
		Category: "Autoscalers",
		Usage:    "Remove an autoscaler from an application",
		Flags:    []cli.Flag{appFlag},
		Description: `Remove an autoscaler for a container type of an application

   Example
     scalingo --app my-app autoscalers-remove web`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "autoscalers-remove")
				return
			}

			currentApp := appdetect.CurrentApp(c)
			err := autoscalers.Remove(currentApp, c.Args()[0])
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "autoscalers-remove")
		},
	}
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
