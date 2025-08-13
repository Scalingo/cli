package cmd

import (
	"context"

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
				_ = cli.ShowCommandHelp(ctx, c, args.First())
				return nil
			}
			_ = cli.ShowAppHelp(c)
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "help")
			_ = autocomplete.HelpAutoComplete(c)
		},
	}
)
