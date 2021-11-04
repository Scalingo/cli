package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/config"
)

var (
	ConfigCommand = cli.Command{
		Name:     "config",
		Category: "Global",
		Usage:    "Configure the CLI",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "region", Value: "", Usage: "Configure the default region used by the CLI"},
		},
		Description: `
   Example
     'scalingo config --region agora-fr1'

	 Can also be configured using the environment variable
	   SCALINGO_REGION=agora-fr1`,
		Action: func(c *cli.Context) {
			if c.String("region") != "" {
				err := config.SetRegion(c.String("region"))
				if err != nil {
					errorQuit(err)
				}
			}

			// If no flag are given, display the current config
			if c.String("region") == "" {
				config.Display()
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "config")
		},
	}
)
