package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	forceHTTPSCommand = cli.Command{
		Name:     "force-https",
		Category: "App Management",
		Usage:    "",
		Flags: []cli.Flag{
			appFlag,
			cli.BoolFlag{Name: "enable, e", Usage: "Enable force HTTPS (default)"},
			cli.BoolFlag{Name: "disable, d", Usage: "Disable force HTTPS"},
		},
		Description: `When enabled, this feature will automatically redirect HTTP traffic to HTTPS for all domains associated with this application.

   Example
     scalingo --app my-app force-https --enable
	 `,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) > 1 {
				cli.ShowCommandHelp(c, "force-https")
				return
			}

			enable := true
			if c.IsSet("disable") {
				enable = false
			}

			err := apps.ForceHTTPS(currentApp, enable)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "force-https")
		},
	}

	stickySessionCommand = cli.Command{
		Name:     "sticky-session",
		Category: "App Management",
		Usage:    "",
		Flags: []cli.Flag{
			appFlag,
			cli.BoolFlag{Name: "enable, e", Usage: "Enable sticky session (default)"},
			cli.BoolFlag{Name: "disable, d", Usage: "Disable sticky session"},
		},
		Description: `When enabled, application user sessions will be sticky: they will always return to the same container.

   Example
     scalingo --app my-app sticky-session --enable
	 `,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) > 1 {
				cli.ShowCommandHelp(c, "sticky-session")
				return
			}

			enable := true
			if c.IsSet("disable") {
				enable = false
			}

			err := apps.StickySession(currentApp, enable)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "sticky-session")
		},
	}

	routerLogsCommand = cli.Command{
		Name:     "router-logs",
		Category: "App Management",
		Usage:    "Enable/disable router logs for your application",
		Flags: []cli.Flag{
			appFlag,
			cli.BoolFlag{Name: "enable, e", Usage: "Enable router logs"},
			cli.BoolFlag{Name: "disable, d", Usage: "Disable router logs (default)"},
		},
		Description: `Enable/disable router logs for your application.

   Example
     scalingo --app my-app router-logs --enable
	 `,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) > 1 {
				cli.ShowCommandHelp(c, "router-logs")
				return
			}

			enable := false
			if c.IsSet("enable") {
				enable = true
			}

			err := apps.RouterLogs(currentApp, enable)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "router-logs")
		},
	}

)
