package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/stacks"
	"github.com/Scalingo/cli/utils"
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

		Action: func(ctx context.Context, c *cli.Command) error {
			err := stacks.List(c.Context, c.Bool("with-deprecated"))
			if err != nil {
				errorQuit(c.Context, err)
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
			Examples:    []string{"scalingo --app my-app stacks-set scalingo-22"},
			SeeAlso:     []string{"stacks"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(ctx, c, "stacks-set")
				return nil
			}
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)

			err := stacks.Set(c.Context, currentApp, c.Args().First())
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			autocomplete.StacksSetAutoComplete(c)
		},
	}
)
