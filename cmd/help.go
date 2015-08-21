package cmd

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	HelpCommand = cli.Command{
		Name:    "help",
		Usage:   "Shows a list of commands or help for one command",
		Aliases: []string{"h"},
		Action: func(c *cli.Context) {
			args := c.Args()
			if args.Present() {
				cli.ShowCommandHelp(c, args.First())
			} else {
				cli.ShowAppHelp(c)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.HelpAutoComplete(c)
		},
	}
)
