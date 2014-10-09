package cmd

import (
	"github.com/Scalingo/cli/addon_resources"
	"github.com/Scalingo/cli/appdetect"
	"github.com/codegangsta/cli"
)

var (
	AddonResourcesListCommand = cli.Command{
		Name:        "addons",
		Description: "List all addons used by your app.",
		Usage:       "List used addons",
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			if err := addon_resources.List(currentApp); err != nil {
				errorQuit(err)
			}
		},
	}
)
