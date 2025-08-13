package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	HelpCommand = cli.Command{
		Name:  "help",
		Usage: "Shows a list of commands or help for one command",
		Action: func(ctx context.Context, c *cli.Command) error {
			args := c.Args()
			if args.Present() {
				cli.ShowCommandHelp(c, args.First())
				return nil
			}
			cli.ShowAppHelp(c)
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			autocomplete.CmdFlagsAutoComplete(c, "help")
			autocomplete.HelpAutoComplete(c)
		},
	}
)
