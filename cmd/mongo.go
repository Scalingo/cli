package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
)

var (
	MongoConsoleCommand = cli.Command{
		Name:     "mongo-console",
		Aliases:  []string{"mongodb-console"},
		Category: "Databases",
		Usage:    "Run an interactive console with your MongoDB addon",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{Name: "size", Aliases: []string{"s"}, Value: "", Usage: "Size of the container"},
			&cli.StringFlag{Name: "env", Aliases: []string{"e"}, Value: "", Usage: "Environment variable name to use for the connection to the database"},
		},
		Description: CommandDescription{
			Description: `Run an interactive console with your MongoDB addon
The --size flag makes it easy to specify the size of the container executing
the MongoDB console. Each container size has different price and performance.
You can read more about container sizes here:
http://doc.scalingo.com/internals/container-sizes.html`,
			Examples: []string{
				"scalingo --app my-app mongo-console",
				"scalingo --app my-app mongo-console --size L",
				"scalingo --app my-app mongo-console --env MY_MONGO_URL",
			},
			SeeAlso: []string{"redis-console", "mysql-console"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(ctx, c, "mongo-console")
				return nil
			}

			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeDBs)

			err := db.MongoConsole(c.Context, db.MongoConsoleOpts{
				App:          currentApp,
				Size:         c.String("s"),
				VariableName: c.String("e"),
			})
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			autocomplete.CmdFlagsAutoComplete(c, "mongo-console")
		},
	}
)
