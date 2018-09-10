package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/urfave/cli"
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
		Before: AuthenticateHook,
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
)
