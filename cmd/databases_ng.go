package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/dbng"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v8"
)

var (
	databasesListCommand = cli.Command{
		Name:     "databases",
		Category: "Databases",
		Usage:    "List the databases next generation that you own",
		Description: CommandDescription{
			Description: "List all the databases next generation of which you are an owner",
			SeeAlso:     []string{"database-info", "database-create", "database-destroy"},
		}.Render(),
		Action: func(ctx context.Context, _ *cli.Command) error {
			err := dbng.List(ctx)
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "databases")
		},
	}

	databaseShowCommand = cli.Command{
		Name:      "database-info",
		Category:  "Databases",
		Usage:     "View database next generation",
		ArgsUsage: "database-id",
		Description: CommandDescription{
			Description: "View database next generation detailed informations",
			Examples:    []string{"scalingo database-info database_id"},
			SeeAlso:     []string{"databases", "database-create", "database-destroy"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			appID := c.Args().First()
			if appID == "" {
				return cli.ShowCommandHelp(ctx, c, "database-info")
			}

			utils.CheckForConsent(ctx, currentApp)

			err := dbng.Show(ctx, appID)
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "database-info")
		},
	}

	databaseAddCommand = cli.Command{
		Name:      "database-create",
		Category:  "Databases",
		Usage:     "Create a database next generation",
		ArgsUsage: "database-name",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "type", Usage: "Database type", Required: true},
			&cli.StringFlag{Name: "plan", Usage: "Database plan", Required: true},
			&cli.StringFlag{Name: "project", Usage: "Project ID", Required: false},
			&cli.BoolFlag{Name: "wait", Usage: "Wait for creation", Required: false},
		},
		Description: CommandDescription{
			Description: "Create a new database next generation",
			Examples: []string{
				"scalingo database-create --type postgresql --plan postgresql-dbng-starter-2048 my_super_database",
				"scalingo database-create --type postgresql --plan postgresql-dbng-starter-2048 --wait my_super_database",
			},
			SeeAlso: []string{"databases", "database-info", "database-destroy"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			appName := c.Args().First()
			if appName == "" {
				return cli.ShowCommandHelp(ctx, c, "database-create")
			}

			utils.CheckForConsent(ctx, currentApp)

			params := scalingo.DatabaseCreateParams{
				AddonProviderID: c.String("type"),
				PlanID:          c.String("plan"),
				ProjectID:       c.String("project"),
				Name:            appName,
			}
			err := dbng.Add(ctx, params, c.Bool("wait"))
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "database-create")
		},
	}

	databaseDeleteCommand = cli.Command{
		Name:      "database-destroy",
		Category:  "Databases",
		Usage:     "Delete database next generation",
		ArgsUsage: "database-id",
		Description: CommandDescription{
			Description: "Delete database next generation",
			Examples:    []string{"scalingo database-destroy database_id"},
			SeeAlso:     []string{"databases", "database-info", "database-create"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			appID := c.Args().First()
			if appID == "" {
				return cli.ShowCommandHelp(ctx, c, "database-destroy")
			}

			utils.CheckForConsent(ctx, currentApp)

			err := dbng.Delete(ctx, appID)
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "database-destroy")
		},
	}
)
