package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/log_drains"
	"github.com/Scalingo/go-scalingo"
	"github.com/urfave/cli"
)

var (
	logDrainsListCommand = cli.Command{
		Name:     "log-drains",
		Category: "Log drains",
		Flags: []cli.Flag{appFlag,
			addonFlag,
			cli.BoolFlag{Name: "with-addons", Usage: "also list the log drains of all addons"},
		},
		Usage: "List the log drains of an application",
		Description: `List all the log drains of an application:

	Use the parameter: "--addon <addon_uuid>" to list log drains of a specific addon
	Use the parameter: "--with-addons" to list log drains of all addons connected to the application

	Examples:
		$ scalingo --app my-app log-drains
		$ scalingo --app my-app log-drains --addon ad-9be0fc04-bee6-4981-a403-a9ddbee7bd1f
		$ scalingo --app my-app log-drains --with-addons

	# See also commands 'log-drains-add', 'log-drains-remove'`,

		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "log-drains")
				return
			}

			var addonID string
			if c.GlobalString("addon") != "<addon_id>" {
				addonID = c.GlobalString("addon")
			} else if c.String("addon") != "<addon_id>" {
				addonID = c.String("addon")
			}

			err := log_drains.List(currentApp, log_drains.ListAddonOpts{
				WithAddons: c.Bool("with-addons"),
				AddonID:    addonID,
			})
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "log-drains")
		},
	}

	logDrainsAddCommand = cli.Command{
		Name:     "log-drains-add",
		Category: "Log drains",
		Usage:    "Add a log drain to an application",
		Flags: []cli.Flag{appFlag,
			addonFlag,
			cli.BoolFlag{Name: "with-addons", Usage: "also add the log drains to all addons"},
			cli.StringFlag{Name: "type", Usage: "Communication protocol", Required: true},
			cli.StringFlag{Name: "url", Usage: "URL of self hosted ELK"},
			cli.StringFlag{Name: "host", Usage: "Host of logs management service"},
			cli.StringFlag{Name: "port", Usage: "Port of logs management service"},
			cli.StringFlag{Name: "token", Usage: "Used by certain vendor for authentication"},
			cli.StringFlag{Name: "drain-region", Usage: "Used by certain logs management service to identify the region of their servers"},
		},
		Description: `Add a log drain to an application:

	Examples:
		$ scalingo --app my-app log-drains-add --type datadog --token 123456789abcdef --drain-region eu-west-2
		$ scalingo --app my-app log-drains-add --type ovh-graylog --token 123456789abcdef --host tag3.logs.ovh.com
		$ scalingo --app my-app log-drains-add --type logentries --token 123456789abcdef
		$ scalingo --app my-app log-drains-add --type papertrail --host logs2.papertrailapp.com --port 12345
		$ scalingo --app my-app log-drains-add --type syslog --host custom.logstash.com --port 12345
		$ scalingo --app my-app log-drains-add --type elk --url https://my-user:123456789abcdef@logstash-app-name.osc-fr1.scalingo.io

	Add a log drain to an addon:

		Use the parameter: "--addon <addon_uuid>" to your add command to add a log drain to a specific addon
		Use the parameter: "--with-addons" to list log drains of all addons connected to the application.

		Warning: At the moment, only databases addons are able to throw logs.

	Examples:
		$ scalingo --app my-app --addon ad-3c2f8c81-99bd-4667-9791-466799bd4667 log-drains-add --type datadog --token 123456789abcdef --drain-region eu-west-2
		$ scalingo --app my-app --with-addons log-drains-add --type datadog --token 123456789abcdef --drain-region eu-west-2

	# See also commands 'log-drains', 'log-drains-remove'`,

		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)

			var addonID string
			if c.GlobalString("addon") != "<addon_id>" {
				addonID = c.GlobalString("addon")
			} else if c.String("addon") != "<addon_id>" {
				addonID = c.String("addon")
			}

			err := log_drains.Add(currentApp,
				log_drains.AddAddonOpts{
					WithAddons: c.Bool("with-addons"),
					AddonID:    addonID,
				},
				scalingo.LogDrainAddParams{
					Type:        c.String("type"),
					URL:         c.String("url"),
					Host:        c.String("host"),
					Port:        c.String("port"),
					Token:       c.String("token"),
					DrainRegion: c.String("drain-region"),
				})
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "log-drains")
		},
	}

	logDrainsRemoveCommand = cli.Command{
		Name:     "log-drains-remove",
		Category: "Log drains",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Remove a log drain from an application",
		Description: `Remove a log drain from an application:

	$ scalingo --app my-app log-drains-remove syslog://custom.logstash.com:12345

	# See also commands 'log-drains-add', 'log-drains'`,

		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			var err error
			if len(c.Args()) == 1 {
				err = log_drains.Remove(currentApp, c.Args()[0])
			} else {
				cli.ShowCommandHelp(c, "log-drains-remove")
			}

			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "log-drains-remove")
			autocomplete.LogDrainsRemoveAutoComplete(c)
		},
	}
)
