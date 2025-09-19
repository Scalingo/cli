package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-utils/errors/v2"
)

var (
	runCommand = cli.Command{
		Name:      "run",
		Aliases:   []string{"r"},
		Category:  "App Management",
		Usage:     "Run any command for your app",
		ArgsUsage: "command-to-execute",
		Flags: []cli.Flag{&appFlag,
			&cli.BoolFlag{Name: "detached", Aliases: []string{"d"}, Usage: "Run a detached container"},
			&cli.StringFlag{Name: "size", Aliases: []string{"s"}, Value: "", Usage: "Size of the container"},
			&cli.StringFlag{Name: "type", Aliases: []string{"t"}, Value: "", Usage: "Procfile Type"},
			&cli.StringSliceFlag{Name: "env", Aliases: []string{"e"}, Usage: "Environment variables"},
			&cli.StringSliceFlag{Name: "file", Aliases: []string{"f"}, Usage: "Files to upload"},
			&cli.BoolFlag{Name: "silent", Usage: "Do not output anything on stderr"},
		},
		Description: `Run command in current app context, a one-off container will be
   start with your application environment loaded.

   Examples
     scalingo --app rails-app run bundle exec rails console
     scalingo --app rails-app run --detached bundle exec rake long:task
     scalingo --app appname run --size XL bash
     scalingo --app symfony-app run php app/console cache:clear
     scalingo --app test-app run --silent custom/command > localoutput

   The --detached flag let you run a 'detached' one-off container, it means the
   container will be started and you'll get back your terminal immediately. Its
   output will be accessible from the logs of the application (command 'logs')
   You can see if the task is still running with the command 'ps' which will
   display the list of the running containers.

   The --size flag makes it easy to specify the size of the container you want
   to run. Each container size has different price and performance. You can read
   more about container sizes here:
   http://doc.scalingo.com/internals/container-sizes.html

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
   useful to you. The option '--file' has been built for this purpose. You can even
   upload multiple files if you wish. You will be able to find these files in the
   '/tmp/uploads' directory of the one-off container. Each file size cannot exceed 100 MiB.

   Example
     scalingo run --file mysqldump.sql rails dbconsole < /tmp/uploads/mysqldump.sql`,
		Action: func(ctx context.Context, c *cli.Command) error {
			currentResource := detect.GetCurrentResource(ctx, c)
			isDB, err := utils.IsResourceDatabase(ctx, currentResource)
			if err != nil && !errors.Is(err, utils.ErrResourceNotFound) {
				errorQuit(ctx, err)
			}
			if isDB {
				io.Error("It is currently impossible to use `run` on a database.")
				return nil
			}

			currentApp := currentResource

			opts := apps.RunOpts{
				App:      currentApp,
				Cmd:      c.Args().Slice(),
				Size:     c.String("s"),
				Type:     c.String("t"),
				CmdEnv:   c.StringSlice("e"),
				Files:    c.StringSlice("f"),
				Silent:   c.Bool("silent"),
				Detached: c.Bool("detached"),
			}
			if (c.Args().Len() == 0 && c.String("t") == "") || (c.Args().Len() > 0 && c.String("t") != "") {
				_ = cli.ShowCommandHelp(ctx, c, "run")
				return nil
			}

			if opts.Detached && len(opts.Files) > 0 {
				io.Error("It is currently impossible to use detached one-off with an uploaded file. Please either remove the --detached or --file flags.")
				return nil
			}

			utils.CheckForConsent(ctx, currentApp)

			err = apps.Run(ctx, opts)
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "run")
		},
	}
)
