package cmd

import (
	"context"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/privatenetworks"
	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/logger"
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
				Value: "json",
				Usage: "[" + outputFormatJSON + "|" + outputFormatTable + "]",
			},
			&cli.StringFlag{
				Name:  "page",
				Value: "1",
				Usage: "[page]",
			},
			&cli.StringFlag{
				Name:  "per-page",
				Value: "20",
				Usage: "[per-page]",
			},
		},
		Usage: "List the private network domain names of an application",
		Description: CommandDescription{
			Description: "List all the private network domain names of an application",
			Examples:    []string{"scalingo --app my-app private-networks-domain-names --format table"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				err := cli.ShowCommandHelp(ctx, c, "private-networks-domain-names")
				if err != nil {
					errorQuit(ctx, err)
				}
				return nil
			}

			pageStr := c.String("page")
			perPageStr := c.String("per-page")
			formatStr := c.String("format")
			ctx, _ = logger.WithFieldsToCtx(ctx, logrus.Fields{
				"page":     pageStr,
				"per_page": perPageStr,
				"format":   formatStr,
			})

			var err error
			var page int
			if pageStr != "" {
				page, err = strconv.Atoi(pageStr)
				if err != nil || page < 1 {
					return errors.New(ctx, "invalid page number")
				}
			}

			var perPage int
			if perPageStr != "" {
				perPage, err = strconv.Atoi(perPageStr)
				if err != nil || perPage < 1 || perPage > 50 {
					return errors.New(ctx, "invalid per_page number")
				}
			}

			currentApp := detect.CurrentApp(c)

			err = privatenetworks.List(ctx, currentApp, formatStr, uint(page), uint(perPage))
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
