package cmd

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/env"
)

var (
	envCommand = cli.Command{
		Name:     "env",
		Category: "Environment",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Display the environment of your apps",
		Description: `List all the environment variables:

    $ scalingo --app my-app env

    # See also commands 'env-get', 'env-set' and 'env-unset'`,

		Action: func(c *cli.Context) {
			currentApp := detect.CurrentApp(c)
			var err error
			if len(c.Args()) == 0 {
				err = env.Display(currentApp)
			} else {
				cli.ShowCommandHelp(c, "env")
			}

			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "env")
		},
	}

	envGetCommand = cli.Command{
		Name:     "env-get",
		Category: "Environment",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Get the requested environment variable from your app",
		Description: `Get the value of a specific environment variable:

    $ scalingo --app my-app env-get VAR1

    # See also commands 'env', 'env-set' and 'env-unset'`,

		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "env")
				return
			}

			currentApp := detect.CurrentApp(c)
			variableValue, err := env.Get(currentApp, c.Args()[0])
			if err != nil {
				errorQuit(err)
			}
			fmt.Println(variableValue)
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "env")
		},
	}

	envSetCommand = cli.Command{
		Name:     "env-set",
		Category: "Environment",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Set the environment variables of your apps",
		Description: `Set variables:

    $ scalingo --app my-app env-set VAR1=VAL1 VAR2=VAL2

    # See also commands 'env', 'env-get' and 'env-unset'`,

		Action: func(c *cli.Context) {
			currentApp := detect.CurrentApp(c)
			var err error
			if len(c.Args()) > 0 {
				err = env.Add(currentApp, c.Args())
			} else {
				cli.ShowCommandHelp(c, "env-set")
				return
			}
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "env-set")
		},
	}

	envUnsetCommand = cli.Command{
		Name:     "env-unset",
		Category: "Environment",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Unset environment variables of your apps",
		Description: `Unset variables:

    $ scalingo --app my-app env-unset VAR1 VAR2

    # See also commands 'env', 'env-get' and 'env-set'`,

		Action: func(c *cli.Context) {
			currentApp := detect.CurrentApp(c)
			var err error
			if len(c.Args()) > 0 {
				err = env.Delete(currentApp, c.Args())
			} else {
				cli.ShowCommandHelp(c, "env-unset")
			}
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "env-unset")
			autocomplete.EnvUnsetAutoComplete(c)
		},
	}
)
