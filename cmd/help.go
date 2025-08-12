package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	HelpCommand = cli.Command{
		Name:  "help",
		Usage: "Shows a list of commands or help for one command",
		Action: func(c *cli.Context) error {
			args := c.Args()
			if args.Present() {
				cli.ShowCommandHelp(c, args.First())
				return nil
			}
			cli.ShowAppHelp(c)
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "help")
			autocomplete.HelpAutoComplete(c)
		},
	}
)
