package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/go-utils/errors"
)

var (
	sendSignalCommand = cli.Command{
		Name:     "send-signal",
		Category: "App Management",
		Usage:    "Send SIGUSR1/2 to your application containers",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{Name: "signal", Aliases: []string{"s"}, Usage: "signal to send to the container"},
		},
		Description: `Send SIGUSR1/2 to your application containers
	Example
	  'scalingo --app my-app send-signal --signal SIGUSR1 web-1'
	  'scalingo --app my-app send-signal --signal SIGUSR2 web-1 web-2'`,
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() == 0 {
				cli.ShowCommandHelp(c, "send-signal")
				return nil
			}

			err := apps.SendSignal(c.Context, currentApp, c.String("signal"), c.Args().Slice())
			if err != nil {
				rootError := errors.RootCause(err)
				errorQuit(rootError)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "send-signal")
		},
	}
)
