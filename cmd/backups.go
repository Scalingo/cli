package cmd

import (
	"fmt"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/db"
	"github.com/urfave/cli"
)

var (
	BackupListCommand = cli.Command{
		Name:     "backups",
		Category: "Addons",
		Usage:    "List backups for an addon",
		Flags:    []cli.Flag{appFlag, addonFlag},
		Description: `  List all available backups for an addon:
		$ scalingo -a myapp --addon my_addon backups

		# See also 'addons' and 'backup-download'
		`,
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			addon := addonName(c)
			err := db.ListBackups(currentApp, addon)
			if err != nil {
				errorQuit(err)
			}
		},
	}

	BackupDownloadCommand = cli.Command{
		Name:     "backup-download",
		Category: "Addons",
		Usage:    "Download a backup",
		Flags: []cli.Flag{appFlag, addonFlag, cli.StringFlag{
			Name:  "backup, b",
			Usage: "Name of the current backup",
		}, cli.StringFlag{
			Name:  "output, o",
			Usage: "Output file (- for stdout)",
		}, cli.BoolFlag{
			Name:  "silent, s",
			Usage: "Do not show meter and loading messages",
		}},
		Description: `  Download a specific backup:
		$ scalingo -a myapp --addon my_addon backup-download --backup my_backup

		# See also 'backups' and 'addons'`,
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			addon := addonName(c)
			backup := c.String("backup")
			if backup == "" {
				fmt.Println("Please specify an backup using the --backup flag")
				return
			}

			opts := db.DownloadBackupOpts{
				Output: c.String("output"),
				Silent: c.Bool("silent"),
			}

			err := db.DownloadBackup(currentApp, addon, backup, opts)
			if err != nil {
				errorQuit(err)
			}
		},
	}
)
