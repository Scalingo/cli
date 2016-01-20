package cmd

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	ScaleCommand = cli.Command{
		Name:      "scale",
		ShortName: "s",
		Category:  "App Management",
		Flags:     []cli.Flag{appFlag, cli.BoolFlag{Name: "synchronous", Usage: "Do the scaling synchronously", EnvVar: ""}},
		Usage:     "Scale your application instantly",
		Description: `Scale your application processes.
   Example
     'scalingo --app my-app scale web:2 worker:1'
     'scalingo --app my-app scale web:1 worker:0'
     'scalingo --app my-app scale web:1:XL'
     'scalingo --app my-app scale web:+1 worker:-1'`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) == 0 {
				cli.ShowCommandHelp(c, "scale")
			} else if err := apps.Scale(currentApp, c.Bool("synchronous"), c.Args()); err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "scale")
			autocomplete.ScaleAutoComplete(c)
		},
	}
)
