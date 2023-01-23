package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/stacks"
)

var (
	stacksListCommand = cli.Command{
		Name:     "stacks",
		Category: "Runtime Stacks",
		Usage:    "List the available runtime stacks",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "with-deprecated", Usage: "Show also deprecated stacks"},
		},
		Description: CommandDescription{
			Description: "List all the available runtime stacks for your apps",
			SeeAlso:     []string{"stacks-set"},
		}.Render(),

		Action: func(c *cli.Context) error {
			err := stacks.List(c.Context, c.Bool("with-deprecated"))
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
	}

	stacksSetCommand = cli.Command{
		Name:      "stacks-set",
		Category:  "Runtime Stacks",
		Usage:     "Set the runtime stack of an app",
		ArgsUsage: "stack-id",
		Flags:     []cli.Flag{&appFlag},
		Description: CommandDescription{
			Description: "Set the runtime stack of an app (deployment cache will be reset)",
			Examples:    []string{"scalingo --app my-app stacks-set scalingo-18"},
			SeeAlso:     []string{"stacks"},
		}.Render(),

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "stacks-set")
				return nil
			}

			err := stacks.Set(c.Context, currentApp, c.Args().First())
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.StacksSetAutoComplete(c)
		},
	}
)
