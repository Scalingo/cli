package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/stacks"
	"github.com/urfave/cli"
)

var (
	stacksListCommand = cli.Command{
		Name:     "stacks",
		Category: "Runtime Stacks",
		Usage:    "List the available runtime stacks",
		Description: ` List all the available rutime stacks for your apps:
    $ scalingo stacks

		# See also 'stacks-set'
`,
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			err := stacks.List()
			if err != nil {
				errorQuit(err)
			}
		},
	}

	stacksSetCommand = cli.Command{
		Name:     "stacks-set",
		Category: "Runtime Stacks",
		Usage:    "Set the runtime stack of an app",
		Flags: []cli.Flag{appFlag, cli.StringFlag{
			Name:  "stack",
			Usage: "Stack to use", Value: "<stack id or name>"},
		},
		Description: ` Set the runtime stack of an app (deployment cache will be reseted):
    $ scalingo -a my-app stacks-set

		# See also 'stacks'
`,
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			var err error
			if len(c.Args()) == 0 {
				err = stacks.Set(currentApp, c.String("stack"))
			} else {
				cli.ShowCommandHelp(c, "stacks-set")
			}
			if err != nil {
				errorQuit(err)
			}
		},
	}
)
