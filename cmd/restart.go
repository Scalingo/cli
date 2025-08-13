package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
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

		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)

			if err := apps.Restart(c.Context, currentApp, c.Bool("s"), c.Args().Slice()); err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},

		ShellComplete: func(ctx context.Context, c *cli.Command) {
			autocomplete.CmdFlagsAutoComplete(c, "restart")
			autocomplete.RestartAutoComplete(c)
		},
	}
)
