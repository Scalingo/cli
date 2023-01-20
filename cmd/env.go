package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/env"
	"github.com/Scalingo/cli/utils"
)

var (
	envCommand = cli.Command{
		Name:     "env",
		Category: "Environment",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "Display the environment variables of your apps",
		Description: CommandDescription{
			Description: "List all the environment variables of your app",
			Examples:    []string{"scalingo --app my-app env"},
			SeeAlso:     []string{"env-get", "env-set", "env-unset"},
		}.Render(),

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			var err error
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "env")
				return nil
			}

			utils.CheckForConsent(c.Context, currentApp)

			err = env.Display(c.Context, currentApp)
			if err != nil {
				errorQuit(err)
			}

			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "env")
		},
	}

	envGetCommand = cli.Command{
		Name:      "env-get",
		Category:  "Environment",
		Flags:     []cli.Flag{&appFlag},
		Usage:     "Get the requested environment variable from your app",
		ArgsUsage: "variable-name",
		Description: CommandDescription{
			Description: "Get the value of a specific environment variable",
			Examples:    []string{"scalingo --app my-app env-get VAR1"},
			SeeAlso:     []string{"env", "env-set", "env-unset"},
		}.Render(),

		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "env")
				return nil
			}

			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp)

			variableValue, err := env.Get(c.Context, currentApp, c.Args().First())
			if err != nil {
				errorQuit(err)
			}
			fmt.Println(variableValue)
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "env")
		},
	}

	envSetCommand = cli.Command{
		Name:      "env-set",
		Category:  "Environment",
		Flags:     []cli.Flag{&appFlag},
		Usage:     "Set the environment variables of your apps",
		ArgsUsage: "variable-assignment...",
		Description: CommandDescription{
			Description: "Set environment variables for the app",
			Examples:    []string{"scalingo --app my-app env-set VAR1=VAL1 VAR2=VAL2"},
			SeeAlso:     []string{"env", "env-get", "env-unset"},
		}.Render(),

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			var err error
			if c.Args().Len() > 0 {
				err = env.Add(c.Context, currentApp, c.Args().Slice())
			} else {
				cli.ShowCommandHelp(c, "env-set")
				return nil
			}
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "env-set")
		},
	}

	envUnsetCommand = cli.Command{
		Name:      "env-unset",
		Category:  "Environment",
		Flags:     []cli.Flag{&appFlag},
		Usage:     "Unset environment variables of your apps",
		ArgsUsage: "variable-name...",
		Description: CommandDescription{
			Description: "Unset variables",
			Examples:    []string{"scalingo --app my-app env-unset VAR1 VAR2"},
			SeeAlso:     []string{"env", "env-get", "env-set"},
		}.Render(),

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			var err error
			if c.Args().Len() > 0 {
				err = env.Delete(c.Context, currentApp, c.Args().Slice())
			} else {
				cli.ShowCommandHelp(c, "env-unset")
			}
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "env-unset")
			autocomplete.EnvUnsetAutoComplete(c)
		},
	}
)
