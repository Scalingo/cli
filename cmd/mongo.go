package cmd

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/db"
)

var (
	MongoConsoleCommand = cli.Command{
		Name:     "mongo-console",
		Category: "Databases",
		Usage:    "Run an interactive console with your MongoDB addon",
		Flags:    []cli.Flag{appFlag},
		Description: ` Run an interactive console with your MongoDB addon.
    $ scalingo -a myapp mongo-console

    # See also 'redis-console' and 'mysql-console'
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "mongo-console")
			} else if err := db.MongoConsole(currentApp); err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "mongo-console")
		},
	}
)
