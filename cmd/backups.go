package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/db"
	"github.com/urfave/cli"
)

var (
	BackupListCommand = cli.Command{
		Name:     "backups",
		Category: "Addons",
		Usage:    "List backups for an addon",
		Flags:    []cli.Flag{appFlag},
		Description: `  List all available backups for an addon:
		$ scalingo -a myapp backups ADDON_ID

		# See also 'addons' and 'backup-download'
		`,
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			var err error
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "backups")
			} else {
				err = db.ListBackups(currentApp, c.Args().First())
			}
			if err != nil {
				errorQuit(err)
			}
		},
	}
)
