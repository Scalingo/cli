package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/logdrains"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v6"
)

var (
	logDrainsListCommand = cli.Command{
		Name:     "log-drains",
		Category: "Log drains",
		Flags: []cli.Flag{&appFlag,
			&addonFlag,
			&cli.BoolFlag{Name: "with-addons", Usage: "also list the log drains of all addons"},
		},
		Usage: "List the log drains of an application",
		Description: CommandDescription{
			Description: `List all the log drains of an application

Use the parameter "--addon <addon_uuid>" to list log drains of a specific addon
Use the parameter "--with-addons" to list log drains of all addons connected to the application`,
			Examples: []string{
				"scalingo --app my-app log-drains",
				"scalingo --app my-app log-drains --addon ad-9be0fc04-bee6-4981-a403-a9ddbee7bd1f",
				"scalingo --app my-app log-drains --with-addons",
			},
			SeeAlso: []string{"log-drains-add", "log-drains-remove"},
		}.Render(),
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "log-drains")
				return nil
			}

			utils.CheckForConsent(c.Context, currentApp)

			addonID := addonUUIDFromFlags(c, currentApp)

			err := logdrains.List(c.Context, currentApp, logdrains.ListAddonOpts{
				WithAddons: c.Bool("with-addons"),
				AddonID:    addonID,
			})
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "log-drains")
		},
	}

	logDrainsAddCommand = cli.Command{
		Name:     "log-drains-add",
		Category: "Log drains",
		Usage:    "Add a log drain to an application",
		Flags: []cli.Flag{&appFlag,
			&addonFlag,
			&cli.BoolFlag{Name: "with-addons", Usage: "also add the log drains to all addons"},
			&cli.BoolFlag{Name: "with-databases", Usage: "also add the log drains to all databases"},
			&cli.StringFlag{Name: "type", Usage: "Communication protocol", Required: true},
			&cli.StringFlag{Name: "url", Usage: "URL of self hosted ELK"},
			&cli.StringFlag{Name: "host", Usage: "Host of logs management service"},
			&cli.StringFlag{Name: "port", Usage: "Port of logs management service"},
			&cli.StringFlag{Name: "token", Usage: "Used by certain vendor for authentication"},
			&cli.StringFlag{Name: "drain-region", Usage: "Used by certain logs management service to identify the region of their servers"},
		},
		Description: CommandDescription{
			Description: `Add a log drain to an application and/or its addons.

Use the --app parameter to add a log drain to your application.

To add a log drain to an addon:

   Use the parameter "--addon <addon_uuid>" to your add command to add a log drain to a specific addon
   Use the parameter "--with-addons" to add log drains of all addons connected to the application.
   Use the parameter "--with-databases" to add log drains of all databases connected to the application.

Warning: At the moment, only databases addons are able to forward logs to a drain.
`,
			Examples: []string{
				"scalingo --app my-app log-drains-add --type appsignal --token 123456789abcdef",
				"scalingo --app my-app log-drains-add --type logtail --token 123456789abcdef",
				"scalingo --app my-app log-drains-add --type datadog --token 123456789abcdef --drain-region eu-west-2",
				"scalingo --app my-app log-drains-add --type ovh-graylog --token 123456789abcdef --host tag3.logs.ovh.com",
				"scalingo --app my-app log-drains-add --type papertrail --host logs2.papertrailapp.com --port 12345",
				"scalingo --app my-app log-drains-add --type syslog --host custom.logstash.com --port 12345",
				"scalingo --app my-app log-drains-add --type syslog --token 123456789abcdef --host custom.logstash.com --port 12345",
				"scalingo --app my-app log-drains-add --type elk --url https://my-user:123456789abcdef@logstash-app-name.osc-fr1.scalingo.io",
				"scalingo --app my-app --addon ad-3c2f8c81-99bd-4667-9791-466799bd4667 log-drains-add --type datadog --token 123456789abcdef --drain-region eu-west-2",
				"scalingo --app my-app --with-addons log-drains-add --type datadog --token 123456789abcdef --drain-region eu-west-2",
			},
			SeeAlso: []string{"log-drains", "log-drains-remove"},
		}.Render(),

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)

			addonID := addonUUIDFromFlags(c, currentApp)

			utils.CheckForConsent(c.Context, currentApp)

			if addonID != "" && (c.Bool("with-addons") || c.Bool("with-databases")) {
				cli.ShowCommandHelp(c, "log-drains-add")
				return nil
			}

			if c.Bool("with-addons") {
				fmt.Println("Warning: At the moment, only database addons are able to forward logs to a drain.")
			}

			err := logdrains.Add(c.Context, currentApp, logdrains.AddDrainOpts{
				WithAddons: c.Bool("with-addons") || c.Bool("with-databases"),
				AddonID:    addonID,
				Params: scalingo.LogDrainAddParams{
					Type:        c.String("type"),
					URL:         c.String("url"),
					Host:        c.String("host"),
					Port:        c.String("port"),
					Token:       c.String("token"),
					DrainRegion: c.String("drain-region"),
				},
			})
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "log-drains-add")
		},
	}

	logDrainsRemoveCommand = cli.Command{
		Name:     "log-drains-remove",
		Category: "Log drains",
		Flags: []cli.Flag{
			&appFlag,
			&addonFlag,
			&cli.BoolFlag{Name: "only-app", Usage: "remove the log drains for the application only"},
		},
		Usage:     "Remove a log drain from an application and its associated addons",
		ArgsUsage: "endpoint-url",
		Description: `Remove a log drain from an application and all its addons:

		$ scalingo --app my-app log-drains-remove syslog://custom.logstash.com:12345

	Remove a log drain from a specific addon:
		Use the parameter "--addon <addon_uuid>" to remove a log drain from a specific addon

		$ scalingo --app my-app --addon ad-3c2f8c81-99bd-4667-9791-466799bd4667 log-drains-remove syslog://custom.logstash.com:12345

	Remove a log drain only for the application:
		Use the parameter "--only-app" to remove a log drain only from the application

		$ scalingo --app my-app --only-app log-drains-remove syslog://custom.logstash.com:12345

	# See also commands 'log-drains-add', 'log-drains'`,

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "log-drains-remove")
				return nil
			}
			drain := c.Args().First()

			addonID := addonUUIDFromFlags(c, currentApp)

			if addonID != "" && c.Bool("only-app") {
				cli.ShowCommandHelp(c, "log-drains-remove")
				return nil
			}

			utils.CheckForConsent(c.Context, currentApp)

			message := "This operation will delete the log drain " + drain
			if addonID == "" && !c.Bool("only-app") {
				// addons + app
				message += " for the application and all its addons"
			} else if addonID != "" && !c.Bool("only-app") {
				// addon only
				message += " for the addon " + addonID
			} else {
				// app only
				message += " for the application " + currentApp
			}
			result := askContinue(message)
			if !result {
				fmt.Println("Aborted")
				return nil
			}

			err := logdrains.Remove(c.Context, currentApp, logdrains.RemoveAddonOpts{
				AddonID: addonID,
				OnlyApp: c.Bool("only-app"),
				URL:     drain,
			})

			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "log-drains-remove")
			autocomplete.LogDrainsRemoveAutoComplete(c)
		},
	}
)

func askContinue(message string) bool {
	result := false
	prompt := &survey.Confirm{
		Message: message + "\n\tConfirm deletion ?",
	}
	survey.AskOne(prompt, &result, nil)
	return result
}
