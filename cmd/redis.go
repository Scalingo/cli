package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
)

var (
	RedisConsoleCommand = cli.Command{
		Name:     "redis-console",
		Category: "Databases",
		Usage:    "Run an interactive console with your Redis addon",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{Name: "size", Aliases: []string{"s"}, Value: "", Usage: "Size of the container"},
			&cli.StringFlag{Name: "env", Aliases: []string{"e"}, Value: "", Usage: "Environment variable name to use for the connection to the database"},
		},
		Description: CommandDescription{
			Description: `Run an interactive console with your Redis addon.
			
The --size flag makes it easy to specify the size of the container executing
the Redis console. Each container size has different price and performance.
You can read more about container sizes here:
http://doc.scalingo.com/internals/container-sizes.html`,
			Examples: []string{
				"scalingo --app my-app redis-console",
				"scalingo --app my-app redis-console --size L",
				"scalingo --app my-app redis-console --env MY_REDIS_URL",
			},
			SeeAlso: []string{"mongo-console", "mysql-console"},
		}.Render(),

		Action: func(c *cli.Context) error {
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "redis-console")
				return nil
			}
			currentApp := detect.CurrentApp(c)

			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeDBs)

			err := db.RedisConsole(c.Context, db.RedisConsoleOpts{
				App:          currentApp,
				Size:         c.String("s"),
				VariableName: c.String("e"),
			})
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "redis-console")
		},
	}
)
