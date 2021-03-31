package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/cli/io"
	"github.com/urfave/cli"
)

var (
	backupsListCommand = cli.Command{
		Name:     "backups",
		Category: "Addons",
		Usage:    "List backups for an addon",
		Flags:    []cli.Flag{appFlag, addonFlag},
		Description: `  List all available backups for an addon:
		$ scalingo --app my-app --addon addon_uuid backups

		# See also 'addons' and 'backups-download'
		`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			addon := addonName(c)
			err := db.ListBackups(currentApp, addon)
			if err != nil {
				errorQuit(err)
			}
		},
	}

	backupsCreateCommand = cli.Command{
		Name:     "backups-create",
		Category: "Addons",
		Usage:    "Ask for a new backup",
		Flags:    []cli.Flag{appFlag, addonFlag},
		Description: `  Ask for a new backup of your addon:
		$ scalingo --app my-app --addon addon_uuid backups-download --backup my_backup

		# See also 'backups' and 'addons'`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			addon := addonName(c)

			err := db.CreateBackup(currentApp, addon)
			if err != nil {
				errorQuit(err)
			}
		},
	}

	backupsDownloadCommand = cli.Command{
		Name:     "backups-download",
		Category: "Addons",
		Usage:    "Download a backup",
		Flags: []cli.Flag{appFlag, addonFlag, cli.StringFlag{
			Name:  "backup, b",
			Usage: "ID of the backup to download",
		}, cli.StringFlag{
			Name:  "output, o",
			Usage: "Output file (- for stdout)",
		}, cli.BoolFlag{
			Name:  "silent, s",
			Usage: "Do not show progress bar and loading messages",
		}},
		Description: `  Download a specific backup:
		$ scalingo --app my-app --addon addon_uuid backups-download --backup my_backup

		# See also 'backups' and 'addons'`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			addon := addonName(c)
			backup := c.String("backup")
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

	backupDownloadCommand = cli.Command{
		Name:        "backup-download",
		Category:    backupsDownloadCommand.Category,
		Usage:       backupsDownloadCommand.Usage,
		Description: backupsDownloadCommand.Description,
		Flags:       backupsDownloadCommand.Flags,
		Before: func(*cli.Context) error {
			io.Warningf("DEPRECATED: please use backups-download instead of this command\n\n")
			return nil
		},
		Action: backupsDownloadCommand.Action,
	}
)
