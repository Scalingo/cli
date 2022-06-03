package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/stacks"
)

var (
	stacksListCommand = cli.Command{
		Name:     "stacks",
		Category: "Runtime Stacks",
		Usage:    "List the available runtime stacks",
		Description: `List all the available runtime stacks for your apps:

		Example:
			scalingo stacks

		# See also 'stacks-set'
`,
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
		Flags:    []cli.Flag{appFlag},
		Description: `Set the runtime stack of an app (deployment cache will be reset):

		Example:
			scalingo --app my-app stacks-set scalingo-18

		# See also 'stacks'
`,
		Action: func(c *cli.Context) {
			currentApp := detect.CurrentApp(c)
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "stacks-set")
				return
			}

			err := stacks.Set(currentApp, c.Args()[0])
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.StacksSetAutoComplete(c)
		},
	}
)
