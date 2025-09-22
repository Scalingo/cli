package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	privateNetworksInfoCommand = cli.Command{
		Name:     "private-networks-info",
		Category: "Private Networks",
		Flags:    []cli.Flag{&projectFlag},
		Usage:    "Display the private network information",
		Description: CommandDescription{
			Description: "Display various private network information such as the network ID, name, region, etc.",
			Examples:    []string{"scalingo --project my-project private-networks-info"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "apps-info")
		},
	}
)
