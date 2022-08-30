package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/io"
)

var (
	backupsListCommand = cli.Command{
		Name:     "backups",
		Category: "Addons",
		Usage:    "List backups for an addon",
		Flags:    []cli.Flag{&appFlag, &addonFlag},
		Description: `  List all available backups for an addon:
		$ scalingo --app my-app --addon addon_uuid backups

		# See also 'addons' and 'backups-download'
		`,
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			addonName := addonNameFromFlags(c)
			if addonName == "" {
				fmt.Println("Unable to find the addon name, please use --addon flag.")
				os.Exit(1)
			}
			err := db.ListBackups(c.Context, currentApp, addonName)
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
	}

	backupsCreateCommand = cli.Command{
		Name:     "backups-create",
		Category: "Addons",
		Usage:    "Ask for a new backup",
		Flags:    []cli.Flag{&appFlag, &addonFlag},
		Description: `  Ask for a new backup of your addon:
		$ scalingo --app my-app --addon addon_uuid backups-download --backup my_backup

		# See also 'backups' and 'addons'`,
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			addonName := addonNameFromFlags(c)
			if addonName == "" {
				fmt.Println("Unable to find the addon name, please use --addon flag.")
				os.Exit(1)
			}

			err := db.CreateBackup(c.Context, currentApp, addonName)
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
	}

	backupsDownloadCommand = cli.Command{
		Name:     "backups-download",
		Category: "Addons",
		Usage:    "Download a backup",
		Flags: []cli.Flag{&appFlag, &addonFlag, &cli.StringFlag{
			Name:    "backup",
			Aliases: []string{"b"},
			Usage:   "ID of the backup to download",
		}, &cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "Output file (- for stdout)",
		}, &cli.BoolFlag{
			Name:    "silent",
			Aliases: []string{"s"},
			Usage:   "Do not show progress bar and loading messages",
		}},
		Description: `  Download a specific backup:
		$ scalingo --app my-app --addon addon_uuid backups-download --backup my_backup

		# See also 'backups' and 'addons'`,
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			addonName := addonNameFromFlags(c)
			if addonName == "" {
				fmt.Println("Unable to find the addon name, please use --addon flag.")
				os.Exit(1)
			}
			backup := c.String("backup")
			opts := db.DownloadBackupOpts{
				Output: c.String("output"),
				Silent: c.Bool("silent"),
			}

			err := db.DownloadBackup(c.Context, currentApp, addonName, backup, opts)
			if err != nil {
				errorQuit(err)
			}
			return nil
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
