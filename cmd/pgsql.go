package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/codegangsta-cli"
)

var (
	PgSQLConsoleCommand = cli.Command{
		Name:     "pgsql-console",
		Category: "Databases",
		Usage:    "Run an interactive console with your PostgreSQL addon",
		Flags: []cli.Flag{appFlag,
			cli.StringFlag{Name: "size, s", Value: "", Usage: "Size of the container"},
		},
		Description: ` Run an interactive console with your PostgreSQL addon.

   Examples
    scalingo --app myapp pgsql-console
    scalingo --app myapp pgsql-console --size L

   The --size flag makes it easy to specify the size of the container executing 
   the PostgreSQL console. Each container size has different price and performance. 
   You can read more about container sizes here:
   http://doc.scalingo.com/internals/container-sizes.html

    # See also 'mongo-console' and 'mysql-console'
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			opts := db.PgSQLConsoleOpts{
				App:  currentApp,
				Size: c.String("s"),
			}
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "pgsql-console")
			} else if err := db.PgSQLConsole(opts); err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "pgsql-console")
		},
	}
)
