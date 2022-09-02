package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/region_migrations"
	scalingo "github.com/Scalingo/go-scalingo/v4"
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
		Description: `Migrate an app to another region.
	 Example
	   'scalingo --app my-app migration-create --to osc-fr1'
		`,
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

			err := region_migrations.Create(currentApp, c.String("to"), c.String("new-name"))
			if err != nil {
				errorQuit(err)
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
		Description: `Run a migration step:
	 Example
	   'scalingo --app my-app migration-run --prepare migration-id'
		`,
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "migration-run")
				return nil
			}
			var step scalingo.RegionMigrationStep
			migrationID := c.Args().First()
			currentApp := detect.CurrentApp(c)
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

			err := region_migrations.Run(currentApp, migrationID, step)
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
	}

	migrationAbortCommand = cli.Command{
		Name:     "migration-abort",
		Category: "Region migrations",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "Abort a migration",
		Description: `Abort a running migration
   Example
	   'scalingo --app my-app migration-abort migration-id'`,
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "migration-run")
				return nil
			}

			migrationID := c.Args().First()
			currentApp := detect.CurrentApp(c)

			err := region_migrations.Abort(currentApp, migrationID)
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
	}

	migrationListCommand = cli.Command{
		Name:     "migrations",
		Category: "Region migrations",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "List all migrations linked to an app",
		Description: `List all migrations linked to an app
   Example
	   'scalingo --app my-app migrations'
		`,
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)

			err := region_migrations.List(currentApp)
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
	}

	migrationFollowCommand = cli.Command{
		Name:     "migration-follow",
		Category: "Region migrations",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "Follow a running migration",
		Description: `Listen for new events on a migration
   Example
	   'scalingo --app my-app migration-follow migration-id'
		`,
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)

			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "migration-follow")
				return nil
			}

			migrationID := c.Args().First()

			err := region_migrations.Follow(currentApp, migrationID)
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.RegionMigrationsAutoComplete(c)
		},
	}
)
