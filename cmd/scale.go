package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/urfave/cli"
)

var (
	scaleCommand = cli.Command{
		Name:      "scale",
		ShortName: "s",
		Category:  "App Management",
		Flags: []cli.Flag{appFlag,
			cli.BoolFlag{Name: "synchronous, s", Usage: "Do the scaling synchronously", EnvVar: ""},
		},
		Usage: "Scale your application instantly",
		Description: `Scale your application processes.
   Example
     'scalingo --app my-app scale web:2 worker:1'
     'scalingo --app my-app scale web:1 worker:0'
     'scalingo --app my-app scale web:1:XL'
     'scalingo --app my-app scale web:+1 worker:-1'
     `,
		Action: func(c *cli.Context) {
			if len(c.Args()) == 0 {
				err := cli.ShowCommandHelp(c, "scale")
				if err != nil {
					errorQuit(err)
				}
				return
			}

			currentApp := appdetect.CurrentApp(c)

			err := apps.Scale(currentApp, c.Bool("s"), c.Args())
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "scale")
			autocomplete.ScaleAutoComplete(c)
		},
	}
)
