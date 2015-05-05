package cmd

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/db"
)

var (
	PgSQLConsoleCommand = cli.Command{
		Name:     "pgsql-console",
		Category: "Databases",
		Usage:    "Run an interactive console with your PostgreSQL addon",
		Flags:    []cli.Flag{appFlag},
		Description: ` Run an interactive console with your PostgreSQL addon.
    $ scalingo -a myapp pgsql-console

		# See also 'mongo-console' and 'mysql-console'
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "pgsql-console")
			} else if err := db.PgSQLConsole(currentApp); err != nil {
				errorQuit(err)
			}
		},
	}
)
