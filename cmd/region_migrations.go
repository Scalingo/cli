package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/region_migrations"
	"github.com/urfave/cli"
)

var (
	migrationCreateCommand = cli.Command{
		Name:     "migrations-create",
		Category: "Region migrations",
		Flags: []cli.Flag{
			appFlag,
			cli.BoolFlag{Name: "follow, f", Usage: "Follow the deployment progress"},
			cli.StringFlag{Name: "to", Usage: "Select the destination region"},
		},
		Usage: "Migrate an app to another region",
		Description: `Migrate an app to another region.
	 Example
	   'scalingo --app my-app migrations-create --to osc-fr1 --follow'
		`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "migrations-create")
				return
			}

			if c.String("to") == "" {
				cli.ShowCommandHelp(c, "migrations-create")
				return
			}

			err := region_migrations.Create(currentApp, c.String("to"))
			if err != nil {
				errorQuit(err)
			}
		},
	}
)
