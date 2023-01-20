package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/config"
)

var (
	ConfigCommand = cli.Command{
		Name:     "config",
		Category: "Global",
		Usage:    "Configure the CLI",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "region", Value: "", Usage: "Configure the default region used by the CLI"},
		},
		Description: CommandDescription{
			Description: "Configure the CLI.\n\nCan also be configured using the environment variable SCALINGO_REGION",
			Examples:    []string{"scalingo config --region agora-fr1"},
		}.Render(),
		Action: func(c *cli.Context) error {
			regionName := regionNameFromFlags(c)
			if regionName != "" {
				err := config.SetRegion(c.Context, regionName)
				if err != nil {
					errorQuit(err)
				}
			}

			// If no flag are given, display the current config
			if regionName == "" {
				config.Display()
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "config")
		},
	}
)
