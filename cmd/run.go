package cmd

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	EnvFlag    = cli.StringSlice([]string{})
	FilesFlag  = cli.StringSlice([]string{})
	RunCommand = cli.Command{
		Name:     "run",
		Category: "App Management",
		Usage:    "Run any command for your app",
		Flags: []cli.Flag{appFlag,
			cli.StringSliceFlag{Name: "env, e", Value: &EnvFlag, Usage: "Environment variables", EnvVar: ""},
			cli.StringSliceFlag{Name: "file, f", Value: &FilesFlag, Usage: "Files to upload", EnvVar: ""},
		},
		Description: `Run command in current app context, a one-off container will be
   start with your application environment loaded.

   Example
     scalingo --app rails-app run bundle exec rails console
     scalingo --app synfony-app run php app/console cache:clear

   If you need to inject additional environment variables, you can use the flag
   '-e'. You can use it multiple time to define multiple variables. These
   variables will override those defined in your application environment.

   Example
     scalingo run -e VARIABLE=VALUE -e VARIABLE2=OTHER_VALUE rails console

   Furthermore, you may want to upload a file, like a database dump or anything
   useful to you. The option '-f' has been built for this purpose, you can even
   upload multiple files if you wish. You will be able to find these files in the
   '/tmp/uploads' directory of the one-off container.

   Example
     scalingo run -f mysqldump.sql rails dbconsole < /tmp/uploads/mysqldump.sql`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			opts := apps.RunOpts{
				App:    currentApp,
				Cmd:    c.Args(),
				CmdEnv: c.StringSlice("e"),
				Files:  c.StringSlice("f"),
			}
			if len(c.Args()) == 0 {
				cli.ShowCommandHelp(c, "run")
			} else if err := apps.Run(opts); err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "run")
		},
	}
)
