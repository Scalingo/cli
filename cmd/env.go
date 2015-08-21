package cmd

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/env"
)

var (
	EnvCommand = cli.Command{
		Name:     "env",
		Category: "Environment",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Display the environment of your apps",
		Description: `List all the environment variables:

    $ scalingo -a myapp env

    # See also commands 'env-set' and 'env-unset'`,

		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
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
	}

	EnvSetCommand = cli.Command{
		Name:     "env-set",
		Category: "Environment",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Set the environment variables of your apps",
		Description: `Set variables:

    $ scalingo -a myapp env-set VAR1=VAL1 VAR2=VAL2

    # See also commands 'env' and 'env-unset'`,

		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			var err error
			if len(c.Args()) > 0 {
				err = env.Add(currentApp, c.Args())
			} else {
				cli.ShowCommandHelp(c, "env-set")
			}
			if err != nil {
				errorQuit(err)
			}
		},
	}

	EnvUnsetCommand = cli.Command{
		Name:     "env-unset",
		Category: "Environment",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Unset environment variables of your apps",
		Description: `Unset variables:

    $ scalingo -a myapp env-unset VAR1 VAR2

    # See also commands 'env' and 'env-set'`,

		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
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
			autocomplete.EnvUnsetAutoComplete(c)
		},
	}
)
