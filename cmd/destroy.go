package cmd

import (
	"fmt"

	"github.com/Scalingo/cli/apps"
	"github.com/codegangsta/cli"
)

var (
	DestroyCommand = cli.Command{
		Name:        "destroy",
		ShortName:   "d",
		Usage:       "Destroy an app /!\\",
		Description: "Destroy an app /!\\ It is not reversible\n  Example:\n    'scalingo destroy my-app'",
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "destroy")
			} else {
				var validationName string
				appName := c.Args()[0]
				fmt.Printf("/!\\ You're going to delete %s, this operation is irreversible.\nTo confirm type the name of the application: ", appName)
				fmt.Scan(&validationName)
				if validationName == appName {
					err := apps.Destroy(appName)
					if err != nil {
						errorQuit(err)
					}
				} else {
					fmt.Printf("'%s' is not '%s', aborting…\n", validationName, appName)
				}
			}
		},
	}
)
