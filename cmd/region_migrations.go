package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/regionmigrations"
	"github.com/Scalingo/cli/utils"
	scalingo "github.com/Scalingo/go-scalingo/v7"
)

var (
	migrationCreateCommand = cli.Command{
		Name:     "migration-create",
		Category: "Region migrations",
		Flags: []cli.Flag{
			&appFlag,
			&cli.StringFlag{Name: "to", Usage: "Select the destination region"},
			&cli.StringFlag{Name: "new-name", Usage: "Name of the app in the destination region (same as origin by default)"},
		},
		Usage: "Start migrating an app to another region",
		Description: CommandDescription{
			Description: "Migrate an app to another region",
			Examples:    []string{"scalingo --app my-app migration-create --to osc-fr1"},
		}.Render(),
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "migration-create")
				return nil
			}

			if c.String("to") == "" {
				cli.ShowCommandHelp(c, "migration-create")
				return nil
			}

			utils.CheckForConsent(c.Context, currentApp)

			err := regionmigrations.Create(c.Context, currentApp, c.String("to"), c.String("new-name"))
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
	}
	migrationRunCommand = cli.Command{
		Name:     "migration-run",
		Category: "Region migrations",
		Flags: []cli.Flag{
			&appFlag,
			&cli.BoolFlag{Name: "prepare", Usage: "Create an empty canvas on the new region"},
			&cli.BoolFlag{Name: "data", Usage: "Import databases (and their data) to the new region"},
			&cli.BoolFlag{Name: "finalize", Usage: "Stop the old app and start the new one"},
		},
		Usage: "Run a specific migration step",
		Description: CommandDescription{
			Description: "Run a migration step",
			Examples:    []string{"scalingo --app my-app migration-run --prepare migration-id"},
		}.Render(),

		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "migration-run")
				return nil
			}
			var step scalingo.RegionMigrationStep
			migrationID := c.Args().First()
			currentApp := detect.CurrentApp(c)

			utils.CheckForConsent(c.Context, currentApp)

			stepsFound := 0
			if c.Bool("prepare") {
				stepsFound++
				step = scalingo.RegionMigrationStepPrepare
			}
			if c.Bool("data") {
				stepsFound++
				step = scalingo.RegionMigrationStepData
			}
			if c.Bool("finalize") {
				stepsFound++
				step = scalingo.RegionMigrationStepFinalize
			}
			if stepsFound != 1 {
				cli.ShowCommandHelp(c, "migration-run")
				return nil
			}

			err := regionmigrations.Run(c.Context, currentApp, migrationID, step)
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
	}

	migrationAbortCommand = cli.Command{
		Name:      "migration-abort",
		Category:  "Region migrations",
		Flags:     []cli.Flag{&appFlag},
		Usage:     "Abort a migration",
		ArgsUsage: "migration-id",
		Description: CommandDescription{
			Description: "Abort a running migration",
			Examples:    []string{"scalingo --app my-app migration-abort migration-id"},
		}.Render(),
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				return cli.ShowCommandHelp(c, "migration-abort")
			}

			migrationID := c.Args().First()
			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp)

			err := regionmigrations.Abort(c.Context, currentApp, migrationID)
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
	}

	migrationListCommand = cli.Command{
		Name:     "migrations",
		Category: "Region migrations",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "List all migrations linked to an app",
		Description: CommandDescription{
			Description: "List all migrations linked to an app",
			Examples:    []string{"scalingo --app my-app migrations"},
		}.Render(),
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)

			err := regionmigrations.List(c.Context, currentApp)
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
	}

	migrationFollowCommand = cli.Command{
		Name:      "migration-follow",
		Category:  "Region migrations",
		Flags:     []cli.Flag{&appFlag},
		Usage:     "Follow a running migration",
		ArgsUsage: "migration-id",
		Description: CommandDescription{
			Description: "Listen for new events on a migration",
			Examples:    []string{"scalingo --app my-app migration-follow migration-id"},
		}.Render(),
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)

			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "migration-follow")
				return nil
			}

			migrationID := c.Args().First()

			err := regionmigrations.Follow(c.Context, currentApp, migrationID)
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.RegionMigrationsAutoComplete(c)
		},
	}
)
