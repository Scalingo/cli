package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
)

var (
	InfluxDBConsoleCommand = cli.Command{
		Name:     "influxdb-console",
		Aliases:  []string{"influx-console"},
		Category: "Databases",
		Usage:    "Run an interactive console with your InfluxDB addon",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{Name: "size", Aliases: []string{"s"}, Value: "", Usage: "Size of the container"},
			&cli.StringFlag{Name: "env", Aliases: []string{"e"}, Value: "", Usage: "Environment variable name to use for the connection to the database"},
		},
		Description: ` Run an interactive console with your InfluxDB addon.

   Examples
    scalingo --app my-app influxdb-console
    scalingo --app my-app influxdb-console --size L
    scalingo --app my-app influxdb-console --env MY_INFLUXDB_URL

   The --size flag makes it easy to specify the size of the container executing
   the InfluxDB console. Each container size has different price and performance.
   You can read more about container sizes here:
   http://doc.scalingo.com/internals/container-sizes.html

    # See also 'mongo-console' and 'mysql-console'
`,
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "influxdb-console")
				return nil
			}

			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeDBs)

			err := db.InfluxDBConsole(c.Context, db.InfluxDBConsoleOpts{
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
			autocomplete.CmdFlagsAutoComplete(c, "influxdb-console")
		},
	}
)
