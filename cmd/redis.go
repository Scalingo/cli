package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/cli/detect"
)

var (
	RedisConsoleCommand = cli.Command{
		Name:     "redis-console",
		Category: "Databases",
		Usage:    "Run an interactive console with your Redis addon",
		Flags: []cli.Flag{appFlag,
			cli.StringFlag{Name: "size, s", Value: "", Usage: "Size of the container"},
			cli.StringFlag{Name: "env, e", Value: "", Usage: "Environment variable name to use for the connection to the database"},
		},
		Description: ` Run an interactive console with your Redis addon.

   Examples
    scalingo --app my-app redis-console
    scalingo --app my-app redis-console --size L
    scalingo --app my-app redis-console --env MY_REDIS_URL

   The --size flag makes it easy to specify the size of the container executing
   the Redis console. Each container size has different price and performance.
   You can read more about container sizes here:
   http://doc.scalingo.com/internals/container-sizes.html

    # See also 'mongo-console' and 'mysql-console'
`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "redis-console")
				return
			}

			err := db.RedisConsole(db.RedisConsoleOpts{
				App:          detect.CurrentApp(c),
				Size:         c.String("s"),
				VariableName: c.String("e"),
			})
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "redis-console")
		},
	}
)
