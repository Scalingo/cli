package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
)

var (
	backupsListCommand = cli.Command{
		Name:     "backups",
		Category: "Addons",
		Usage:    "List backups for an addon",
		Flags:    []cli.Flag{&appFlag, &addonFlag},
		Description: CommandDescription{
			Description: "List all available backups for an addon",
			Examples:    []string{"scalingo --app my-app --addon addon_uuid backups"},
			SeeAlso:     []string{"addons", "backups-download"},
		}.Render(),

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			addonName := addonNameFromFlags(c, true)

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
		Description: CommandDescription{
			Description: "Ask for a new backup of your addon",
			Examples:    []string{"scalingo --app my-app --addon addon_uuid backups-create"},
			SeeAlso:     []string{"backups", "addons"},
		}.Render(),
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			addonName := addonNameFromFlags(c, true)

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
			Name:        "backup",
			Aliases:     []string{"b"},
			Usage:       "ID of the backup to download",
			DefaultText: "last successful backup",
		}, &cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "Output file (- for stdout)",
		}, &cli.BoolFlag{
			Name:    "silent",
			Aliases: []string{"s"},
			Usage:   "Do not show progress bar and loading messages",
		}},
		Description: CommandDescription{
			Description: "Download a specific backup",
			Examples:    []string{"scalingo --app my-app --addon addon_uuid backups-download --backup my_backup"},
			SeeAlso:     []string{"backups", "addons"},
		}.Render(),
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			addonName := addonNameFromFlags(c, true)

			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeDBs)

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
