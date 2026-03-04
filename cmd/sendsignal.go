package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/go-utils/errors/v3"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
)

var (
	sendSignalCommand = cli.Command{
		Name:      "send-signal",
		Aliases:   []string{"kill"},
		Category:  "App Management",
		Usage:     "Send SIGUSR1 or SIGUSR2 to your application containers",
		ArgsUsage: "container-type...",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{Name: "signal", Aliases: []string{"s"}, Usage: "signal to send to the container", Required: true},
		},
		Description: CommandDescription{
			Description: "Send SIGUSR1 or SIGUSR2 to your application containers",
			Examples: []string{
				"scalingo --app my-app send-signal --signal SIGUSR1 web-1",
				"scalingo --app my-app send-signal --signal SIGUSR2 web-1 web-2",
				"scalingo --app my-app send-signal --signal SIGUSR2 web",
			},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() == 0 {
				err := cli.ShowCommandHelp(ctx, c, "send-signal")
				if err != nil {
					return errors.Wrapf(ctx, err, "fail to show command helper")
				}
				return nil
			}
			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)

			err := apps.SendSignal(ctx, currentApp, c.String("signal"), c.Args().Slice())
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "send-signal")
		},
	}
)
