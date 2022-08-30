package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/cli/detect"
)

var (
	MySQLConsoleCommand = cli.Command{
		Name:     "mysql-console",
		Category: "Databases",
		Usage:    "Run an interactive console with your MySQL addon",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{Name: "size", Aliases: []string{"s"}, Value: "", Usage: "Size of the container"},
			&cli.StringFlag{Name: "env", Aliases: []string{"e"}, Value: "", Usage: "Environment variable name to use for the connection to the database"},
		},
		Description: ` Run an interactive console with your MySQL addon.

   Examples
    scalingo --app my-app mysql-console
    scalingo --app my-app mysql-console --size L
    scalingo --app my-app mysql-console --env MY_MYSQL_URL

   The --size flag makes it easy to specify the size of the container executing
   the MySQL console. Each container size has different price and performance.
   You can read more about container sizes here:
   http://doc.scalingo.com/internals/container-sizes.html

    # See also 'mongo-console' and 'pgsql-console'
`,
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "mysql-console")
				return nil
			}

			err := db.MySQLConsole(c.Context, db.MySQLConsoleOpts{
				App:          detect.CurrentApp(c),
				Size:         c.String("s"),
				VariableName: c.String("e"),
			})
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "mysql-console")
		},
	}
)
