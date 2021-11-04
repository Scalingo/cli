package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/db"
)

var (
	MongoConsoleCommand = cli.Command{
		Name:     "mongo-console",
		Category: "Databases",
		Usage:    "Run an interactive console with your MongoDB addon",
		Flags: []cli.Flag{appFlag,
			cli.StringFlag{Name: "size, s", Value: "", Usage: "Size of the container"},
		},
		Description: ` Run an interactive console with your MongoDB addon.

   Examples
    scalingo --app myapp mongo-console
    scalingo --app myapp mongo-console --size L

   The --size flag makes it easy to specify the size of the container executing
   the MongoDB console. Each container size has different price and performance.
   You can read more about container sizes here:
   http://doc.scalingo.com/internals/container-sizes.html

    # See also 'redis-console' and 'mysql-console'
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			opts := db.MongoConsoleOpts{
				App:  currentApp,
				Size: c.String("s"),
			}
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "mongo-console")
			} else if err := db.MongoConsole(opts); err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "mongo-console")
		},
	}
)
