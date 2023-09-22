package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
)

var (
	scaleCommand = cli.Command{
		Name:     "scale",
		Aliases:  []string{"s"},
		Category: "App Management",
		Flags: []cli.Flag{&appFlag,
			&cli.BoolFlag{Name: "synchronous", Aliases: []string{"s"}, Usage: "Do the scaling synchronously"},
		},
		Usage:     "Scale your application instantly",
		ArgsUsage: "[scaling-instruction...]",
		Description: CommandDescription{
			Description: "Scale your application processes. Without argument, this command lists the container types declared in your application",
			Examples: []string{
				"scalingo --app my-app scale web:2 worker:1",
				"scalingo --app my-app scale web:1 worker:0",
				"scalingo --app my-app scale web:1:XL",
				"scalingo --app my-app scale web:+1 worker:-1",
			},
		}.Render(),
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)

			if c.Args().Len() == 0 {
				err := apps.ContainerTypes(c.Context, currentApp)
				if err != nil {
					errorQuit(c.Context, err)
				}
				return nil
			}

			err := apps.Scale(c.Context, currentApp, c.Bool("s"), c.Args().Slice())
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "scale")
			autocomplete.ScaleAutoComplete(c)
		},
	}
)
