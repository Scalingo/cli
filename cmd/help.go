package cmd

import (
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/urfave/cli"
)

var (
	HelpCommand = cli.Command{
		Name:  "help",
		Usage: "Shows a list of commands or help for one command",
		Action: func(c *cli.Context) {
			args := c.Args()
			if args.Present() {
				cli.ShowCommandHelp(c, args.First())
				return
			}
			cli.ShowAppHelp(c)
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "help")
			autocomplete.HelpAutoComplete(c)
		},
	}
)
