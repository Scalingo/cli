package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/detect"
)

var (
	openCommand = cli.Command{
		Name:     "open",
		Category: "App Management",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "Open app on default web browser",
		Description: CommandDescription{
			Description: "Open app on default web browser",
			Examples:    []string{"scalingo --app my-app open"},
		}.Render(),
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "open")
				return nil
			}

			currentApp := detect.CurrentApp(c)
			currentRegion := config.C.ScalingoRegion

			err := apps.Open(currentApp, currentRegion)
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "open")
		},
	}
)
