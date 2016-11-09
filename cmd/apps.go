package cmd

import (
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/codegangsta-cli"
)

var (
	AppsCommand = cli.Command{
		Name:        "apps",
		Category:    "Global",
		Description: "List your apps and give some details about them",
		Usage:       "List your apps",
		Before:      AuthenticateHook,
		Action: func(c *cli.Context) {
			if err := apps.List(); err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "apps")
		},
	}
)
