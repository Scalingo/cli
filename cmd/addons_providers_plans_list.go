package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/addonproviders"
	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	AddonProvidersPlansCommand = cli.Command{
		Name:      "addons-plans",
		Category:  "Addons - Global",
		ArgsUsage: "addon-id",
		Description: CommandDescription{
			Description: "List the plans for an addon",
			Examples:    []string{"scalingo addon-plans scalingo-mongodb"},
		}.Render(),
		Usage: "List plans",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() != 1 {
				_ = cli.ShowCommandHelp(ctx, c, "addons-plans")
				return nil
			}
			if err := addonproviders.Plans(ctx, c.Args().First()); err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			autocomplete.AddonsPlansAutoComplete(c)
		},
	}
)
