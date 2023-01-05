package cmd

import (
	"github.com/urfave/cli/v2"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/go-utils/errors"
)

var (
	sendSignalCommand = cli.Command{
		Name:     "send-signal",
		Aliases:  []string{"kill"},
		Category: "App Management",
		Usage:    "Send SIGUSR1 or SIGUSR2 to your application containers",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{Name: "signal", Aliases: []string{"s"}, Usage: "signal to send to the container", Required: true},
		},
		Description: `Send SIGUSR1 or SIGUSR2 to your application containers
	Example
	  'scalingo --app my-app send-signal --signal SIGUSR1 web-1'
	  'scalingo --app my-app send-signal --signal SIGUSR2 web-1 web-2'
	  'scalingo --app my-app send-signal --signal SIGUSR2 web'`,
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() == 0 {
				err := cli.ShowCommandHelp(c, "send-signal")
				if err != nil {
					return errgo.Notef(err, "fail to show command helper")
				}
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
			_ = autocomplete.CmdFlagsAutoComplete(c, "send-signal")
		},
	}
)
