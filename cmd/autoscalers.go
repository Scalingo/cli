package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/autoscalers"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v6"
)

var (
	autoscalersListCommand = cli.Command{
		Name:        "autoscalers",
		Category:    "Autoscalers",
		Usage:       "List the autoscalers of an application",
		Flags:       []cli.Flag{&appFlag},
		Description: "List all the autoscalers of an application and display information about them.",
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "autoscalers")
				return nil
			}
			currentApp := detect.CurrentApp(c)

			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)

			err := autoscalers.List(c.Context, detect.CurrentApp(c))
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "autoscalers")
		},
	}

	autoscalersAddCommand = cli.Command{
		Name:     "autoscalers-add",
		Category: "Autoscalers",
		Usage:    "Add an autoscaler to an application",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{Name: "container-type", Aliases: []string{"c"}, Usage: "Specify the container type affected by the autoscaler"},
			&cli.StringFlag{Name: "metric", Aliases: []string{"m"}, Usage: "Specify the metric you want the autoscaling to apply on"},
			&cli.Float64Flag{Name: "target", Aliases: []string{"t"}, Usage: "Target value for the metric the autoscaler will maintain"},
			&cli.IntFlag{Name: "min-containers", Usage: "lower limit the autoscaler will never scale below"},
			&cli.IntFlag{Name: "max-containers", Usage: "upper limit the autoscaler will never scale above"},
		},
		Description: CommandDescription{
			Description: "Add an autoscaler to an application. It will automatically scale the application up or down depending on the target defined for the given metric.\n\nAll options are mandatory.",
			Examples:    []string{"scalingo --app my-app autoscalers-add --container-type web --metric cpu --target 0.75 --min-containers 2 --max-containers 4"},
		}.Render(),

		Action: func(c *cli.Context) error {
			if !isValidAutoscalerAddOpts(c) {
				err := cli.ShowCommandHelp(c, "autoscalers-add")
				if err != nil {
					errorQuit(c.Context, err)
				}
				return nil
			}

			currentApp := detect.CurrentApp(c)

			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)

			err := autoscalers.Add(c.Context, currentApp, scalingo.AutoscalerAddParams{
				ContainerType: c.String("c"),
				Metric:        c.String("m"),
				Target:        c.Float64("t"),
				MinContainers: c.Int("min-containers"),
				MaxContainers: c.Int("max-containers"),
			})
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "autoscalers-add")
		},
	}

	autoscalersUpdateCommand = cli.Command{
		Name:     "autoscalers-update",
		Category: "Autoscalers",
		Usage:    "Update an autoscaler",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{Name: "container-type", Aliases: []string{"c"}, Usage: "Specify the container type affected by the autoscaler"},
			&cli.StringFlag{Name: "metric", Aliases: []string{"m"}, Usage: "Specify the metric you want the autoscaling to apply on"},
			&cli.Float64Flag{Name: "target", Aliases: []string{"t"}, Usage: "Target value for the metric the autoscaler will maintain"},
			&cli.IntFlag{Name: "min-containers", Usage: "lower limit the autoscaler will never scale below"},
			&cli.IntFlag{Name: "max-containers", Usage: "upper limit the autoscaler will never scale above"},
			&cli.BoolFlag{Name: "disabled", Aliases: []string{"d"}, Usage: "disable/enable the given autoscaler"},
		},
		Description: CommandDescription{
			Description: "Update an autoscaler.\n\nThe 'container-type' option is mandatory.",
			Examples: []string{
				"scalingo --app my-app autoscalers-update --container-type web --max-containers 5",
				"scalingo --app my-app autoscalers-update --container-type web --metric p95_response_time --target 67",
			},
		}.Render(),

		Action: func(c *cli.Context) error {
			if c.Args().Len() != 0 || !c.IsSet("c") {
				err := cli.ShowCommandHelp(c, "autoscalers-update")
				if err != nil {
					errorQuit(c.Context, err)
				}
				return nil
			}

			currentApp := detect.CurrentApp(c)

			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)

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
			err := autoscalers.Update(c.Context, currentApp, c.String("c"), params)
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "autoscalers-update")
		},
	}

	autoscalersEnableCommand = cli.Command{
		Name:      "autoscalers-enable",
		Category:  "Autoscalers",
		Usage:     "Enable an autoscaler",
		ArgsUsage: "container-type",
		Flags:     []cli.Flag{&appFlag},
		Description: CommandDescription{
			Description: "Enable an autoscaler",
			Examples:    []string{"scalingo --app my-app autoscalers-enable web"},
		}.Render(),
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				err := cli.ShowCommandHelp(c, "autoscalers-enable")
				if err != nil {
					errorQuit(c.Context, err)
				}
				return nil
			}

			currentApp := detect.CurrentApp(c)

			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)

			disabled := false
			err := autoscalers.Update(c.Context, currentApp, c.Args().First(), scalingo.AutoscalerUpdateParams{
				Disabled: &disabled,
			})
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "autoscalers-enable")
		},
	}

	autoscalersDisableCommand = cli.Command{
		Name:      "autoscalers-disable",
		Category:  "Autoscalers",
		Usage:     "Disable an autoscaler",
		ArgsUsage: "container-type",
		Flags:     []cli.Flag{&appFlag},
		Description: CommandDescription{
			Description: "Disable an autoscaler",
			Examples:    []string{"scalingo --app my-app autoscalers-disable web"},
		}.Render(),

		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				err := cli.ShowCommandHelp(c, "autoscalers-disable")
				if err != nil {
					errorQuit(c.Context, err)
				}
				return nil
			}

			currentApp := detect.CurrentApp(c)

			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)

			disabled := true
			err := autoscalers.Update(c.Context, currentApp, c.Args().First(), scalingo.AutoscalerUpdateParams{
				Disabled: &disabled,
			})
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "autoscalers-disable")
		},
	}

	autoscalersRemoveCommand = cli.Command{
		Name:      "autoscalers-remove",
		Category:  "Autoscalers",
		Usage:     "Remove an autoscaler from an application",
		ArgsUsage: "container-type",
		Flags:     []cli.Flag{&appFlag},
		Description: CommandDescription{
			Description: "Remove an autoscaler for a container type of an application",
			Examples:    []string{"scalingo --app my-app autoscalers-remove web"},
		}.Render(),

		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "autoscalers-remove")
				return nil
			}

			currentApp := detect.CurrentApp(c)

			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)

			err := autoscalers.Remove(c.Context, currentApp, c.Args().First())
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "autoscalers-remove")
		},
	}
)

func isValidAutoscalerAddOpts(c *cli.Context) bool {
	if c.Args().Len() > 0 {
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
