package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/autoscalers"
	"github.com/Scalingo/cli/cmd/autocomplete"
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

	/*
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
			Description: "Invite someone to collaborate on an application, an invitation will be sent to the given email\n scalingo -a myapp autoscalers-add user@example.com",
			Before:      AuthenticateHook,
			Action: func(c *cli.Context) {
				currentApp := appdetect.CurrentApp(c)
				if len(c.Args()) != 1 {
					cli.ShowCommandHelp(c, "autoscalers-add")
				} else {
					err := autoscalers.Add(currentApp, c.Args()[0])
					if err != nil {
						errorQuit(err)
					}
				}
			},
			BashComplete: func(c *cli.Context) {
				autocomplete.CmdFlagsAutoComplete(c, "autoscalers-add")
				autocomplete.AutoscalersAddAutoComplete(c)
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
