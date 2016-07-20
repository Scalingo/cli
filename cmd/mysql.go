package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/codegangsta-cli"
)

var (
	MySQLConsoleCommand = cli.Command{
		Name:     "mysql-console",
		Category: "Databases",
		Usage:    "Run an interactive console with your MySQL addon",
		Flags:    []cli.Flag{appFlag},
		Description: ` Run an interactive console with your MySQL addon.
    $ scalingo -a myapp mysql-console

    # See also 'mongo-console' and 'pgsql-console'
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "mysql-console")
			} else if err := db.MySQLConsole(currentApp); err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "mysql-console")
		},
	}
)
