package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/cli/detect"
)

var (
	MySQLConsoleCommand = cli.Command{
		Name:     "mysql-console",
		Category: "Databases",
		Usage:    "Run an interactive console with your MySQL addon",
		Flags: []cli.Flag{appFlag,
			cli.StringFlag{Name: "size, s", Value: "", Usage: "Size of the container"},
		},
		Description: ` Run an interactive console with your MySQL addon.

   Examples
    scalingo --app my-app mysql-console
    scalingo --app my-app mysql-console --size L

   The --size flag makes it easy to specify the size of the container executing
   the MySQL console. Each container size has different price and performance.
   You can read more about container sizes here:
   http://doc.scalingo.com/internals/container-sizes.html

    # See also 'mongo-console' and 'pgsql-console'
`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "mysql-console")
				return
			}

			err := db.MySQLConsole(db.MySQLConsoleOpts{
				App:  detect.CurrentApp(c),
				Size: c.String("s"),
			})
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "mysql-console")
		},
	}
)
