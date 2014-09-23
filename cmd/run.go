package cmd

import (
	"github.com/Appsdeck/appsdeck/appdetect"
	"github.com/Appsdeck/appsdeck/apps"
	"github.com/Appsdeck/appsdeck/auth"
	"github.com/codegangsta/cli"
)

var (
	EnvFlag    = cli.StringSlice([]string{})
	FilesFlag  = cli.StringSlice([]string{})
	RunCommand = cli.Command{
		Name:      "run",
		ShortName: "r",
		Usage:     "Run any command for your app",
		Flags: []cli.Flag{
			cli.StringSliceFlag{Name: "env, e", Value: &EnvFlag, Usage: "Environment variables", EnvVar: ""},
			cli.StringSliceFlag{Name: "file, f", Value: &FilesFlag, Usage: "Files to upload", EnvVar: ""},
		},
		Description: `Run command in current app context, your application
   environment will be loaded and you can execute any task.
     Example
       'appsdeck --app my-app run bundle exec rails console'
       'appsdeck --app synfony-app run php app/console cache:clear --env=prod'`,
		Action: func(c *cli.Context) {
			auth.InitAuth()
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			if len(c.Args()) == 0 {
				cli.ShowCommandHelp(c, "run")
			} else if err := apps.Run(currentApp, c.Args(), c.StringSlice("e"), c.StringSlice("f")); err != nil {
				errorQuit(err)
			}
		},
	}
)
