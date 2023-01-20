package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
)

var (
	RestartCommand = cli.Command{
		Name:     "restart",
		Category: "App Management",
		Usage:    "Restart processes of your app",
		Flags: []cli.Flag{&appFlag,
			&cli.BoolFlag{Name: "synchronous", Aliases: []string{"s"}, Usage: "Do the restart synchronously"},
		},
		Description: CommandDescription{
			Description: "Restart one or several process or your application",
			Examples: []string{
				"scalingo --app my-app restart        # Restart all the processes",
				"scalingo --app my-app restart web    # Restart all the web processes",
				"scalingo --app my-app restart web-1  # Restart a specific container",
			},
		}.Render(),

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if err := apps.Restart(c.Context, currentApp, c.Bool("s"), c.Args().Slice()); err != nil {
				errorQuit(err)
			}
			return nil
		},

		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "restart")
			autocomplete.RestartAutoComplete(c)
		},
	}
)
