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
		Flags:    []cli.Flag{appFlag},
		Usage:    "List the log drains of an application",
		Description: `List all the log drains of an application:

	$ scalingo --app my-app log-drains

	# See also commands 'log-drains-add'`,

		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "log-drains")
				return
			}

			err := log_drains.List(currentApp)
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

	# See also commands 'log-drains'`,

		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)

			err := log_drains.Add(currentApp, scalingo.LogDrainAddParams{
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
)
