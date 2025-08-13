package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/addonproviders"
	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	AddonProvidersListCommand = cli.Command{
		Name:        "addons-list",
		Category:    "Addons - Global",
		Description: "List all addons you can add to your app",
		Usage:       "List all addons",
		Action: func(ctx context.Context, c *cli.Command) error {
			if err := addonproviders.List(ctx); err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "addons-list")
		},
	}
)
