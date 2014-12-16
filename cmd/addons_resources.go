package cmd

import (
	"github.com/Scalingo/cli/addon_resources"
	"github.com/Scalingo/cli/appdetect"
	"github.com/codegangsta/cli"
)

var (
	AddonResourcesListCommand = cli.Command{
		Name:     "addons",
		Category: "Addons",
		Usage:    "List used addons",
		Description: ` List all addons used by your app:
    $ scalingo -a myapp addons

    Provision a new addon for an app:
    $ scalingo -a myapp addons <addon-name> provision <plan>

    Destroy an addon of an app:
    $ scalingo -a myapp addons <addon-id> destroy
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			var err error
			if len(c.Args()) == 0 {
				err = addon_resources.List(currentApp)
			} else if len(c.Args()) == 2 && c.Args().Get(1) == "destroy" {
				err = addon_resources.Destroy(currentApp, c.Args().Get(0))
			} else if len(c.Args()) == 3 && c.Args().Get(1) == "upgrade" {
				err = addon_resources.Upgrade(currentApp, c.Args().Get(0), c.Args().Get(2))
			} else if len(c.Args()) == 3 && c.Args().Get(1) == "provision" {
				err = addon_resources.Provision(currentApp, c.Args().Get(0), c.Args().Get(2))
			} else {
				cli.ShowCommandHelp(c, "addons")
			}

			if err != nil {
				errorQuit(err)
			}
		},
	}
)
