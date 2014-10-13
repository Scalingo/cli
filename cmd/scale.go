package cmd

import "github.com/codegangsta/cli"

// import (
// 	"github.com/Scalingo/cli/appdetect"
// 	"github.com/Scalingo/cli/apps"
// 	"github.com/Scalingo/cli/auth"
// 	"github.com/codegangsta/cli"
// )

var (
	// flag         = cli.StringSlice([]string{})
	ScaleCommand = cli.Command{
		Name:      "scale",
		ShortName: "s",
	}

// 	Usage:     "Run any command for your app",
// 	Flags: []cli.Flag{
// 		cli.StringSliceFlag{"env, e", &flag, "Environment variables", ""},
// 	},
// 	Description: `Run command in current app context, your application
// environment will be loaded and you can execute any task.
// Example
// 'appsdeck --app my-app run bundle exec rails console'
// 'appsdeck --app synfony-app run php app/console cache:clear --env=prod'`,
// 	Action: func(c *cli.Context) {
// 		currentApp := appdetect.CurrentApp(c.GlobalString("app"))
// 		if len(c.Args()) == 0 {
// 			cli.ShowCommandHelp(c, "run")
// 		} else if err := apps.Run(currentApp, c.Args(), c.StringSlice("e")); err != nil {
// 			errorQuit(err)
// 		}
// 	},
// }
)
