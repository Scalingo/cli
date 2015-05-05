package cmd

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/db"
)

var (
	RedisConsoleCommand = cli.Command{
		Name:     "redis-console",
		Category: "Databases",
		Usage:    "Run an interactive console with your redis addon",
		Flags:    []cli.Flag{appFlag},
		Description: ` Run an interactive console with your redis addon.
    $ scalingo -a myapp redis-console

    # See also 'mongo-console' and 'mysql-console'
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "redis-console")
			} else if err := db.RedisConsole(currentApp); err != nil {
				errorQuit(err)
			}
		},
	}
)
