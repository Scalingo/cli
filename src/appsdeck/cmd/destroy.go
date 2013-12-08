package cmd

import (
	"appsdeck/apps"
	"appsdeck/auth"
	"fmt"
	"github.com/codegangsta/cli"
)

var (
	DestroyCommand = cli.Command{
		Name:        "destroy",
		ShortName:   "d",
		Description: "/!\\ Destroy an app",
		Usage:       "appsdeck destroy <id or canonical name>",
		Action: func(c *cli.Context) {
			auth.InitAuth()
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "destroy")
			} else {
				var validationName string
				appName := c.Args()[0]
				fmt.Printf("/!\\ Your going to delete %s, this operation is irreversible.\nTo confirm type the name of the application: ", appName)
				fmt.Scan(&validationName)
				if validationName == appName {
					apps.Destroy(appName)
				} else {
					fmt.Printf("'%s' is not '%s', aborting…\n", validationName, appName)
				}
			}
		},
	}
)
