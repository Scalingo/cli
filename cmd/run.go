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
		Name:      "run",
		ShortName: "r",
		Category:  "App Management",
		Usage:     "Run any command for your app",
		Flags: []cli.Flag{appFlag,
			cli.StringFlag{Name: "type, t", Value: "", Usage: "Procfile Type"},
			cli.StringSliceFlag{Name: "env, e", Value: &EnvFlag, Usage: "Environment variables"},
			cli.StringSliceFlag{Name: "file, f", Value: &FilesFlag, Usage: "Files to upload"},
			cli.BoolFlag{Name: "silent, s", Usage: "Do not output anything on stderr"},
		},
		Description: `Run command in current app context, a one-off container will be
   start with your application environment loaded.

   Examples
     scalingo --app rails-app run bundle exec rails console
     scalingo --app synfony-app run php app/console cache:clear
     scalingo --app test-app run --silent custom/command > localoutput

   The --silent flag makes that the only output of the command will be the output
   of the one-off container. There won't be any noise from the command tool itself.

   Thank to the --type flag, you can build shortcuts to commands of your Procfile.
   If your procfile is:

   ==== Procfile
   web: bundle exec rails server
   migrate: bundle rake db:migrate
   ====

   You can run the migrate task with the following command:

   Example:
     scalingo --app my-app run -t migrate

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
				Type:   c.String("t"),
				CmdEnv: c.StringSlice("e"),
				Files:  c.StringSlice("f"),
				Silent: c.Bool("s"),
			}
			if (len(c.Args()) == 0 && c.String("t") == "") || (len(c.Args()) > 0 && c.String("t") != "") {
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
