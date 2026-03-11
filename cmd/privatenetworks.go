package cmd

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/privatenetworks"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/go-utils/pagination"
)

const (
	outputFormatJSON  = "json"
	outputFormatTable = "table"
)

var (
	privateNetworksApplicationDomainsListCommand = cli.Command{
		Name:     "private-networks-domain-names",
		Category: "Private Networks",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{
				Name:  "format",
				Value: outputFormatTable,
				Usage: "[" + outputFormatJSON + "|" + outputFormatTable + "]",
			},
			&cli.IntFlag{
				Name:  "page",
				Value: 1,
				Usage: "[page]",
			},
			&cli.IntFlag{
				Name:  "per-page",
				Value: 20,
				Usage: "[per-page]",
			},
		},
		Usage: "List the private network domain names of an application",
		Description: CommandDescription{
			Description: "List all the private network domain names of an application",
			Examples:    []string{"scalingo --app my-app private-networks-domain-names --format table"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() > 0 {
				err := cli.ShowCommandHelp(ctx, c, "private-networks-domain-names")
				if err != nil {
					errorQuit(ctx, err)
				}
				return nil
			}

			page := c.Int("page")
			perPage := c.Int("per-page")
			formatStr := c.String("format")
			ctx, _ = logger.WithFieldsToCtx(ctx, logrus.Fields{
				"page":     page,
				"per_page": perPage,
				"format":   formatStr,
			})

			currentApp := detect.CurrentApp(ctx, c)

			err := privatenetworks.List(ctx, currentApp, formatStr, pagination.NewRequest(page, perPage))
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "private-networks-domain-names")
		},
	}
)
