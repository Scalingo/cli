package cmd

import (
	"regexp"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
)

var (
	oneOffStopCommand = cli.Command{
		Name:      "one-off-stop",
		Category:  "App Management",
		Usage:     "Stop a running one-off container",
		Flags:     []cli.Flag{&appFlag},
		ArgsUsage: "container-id",
		Description: CommandDescription{
			Description: "Stop a running one-off container",
			Examples: []string{
				"scalingo --app my-app one-off-stop one-off-1234",
				"scalingo --app my-app one-off-stop 1234",
			},
		}.Render(),
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "one-off-stop")
				return nil
			}
			oneOffLabel := c.Args().First()

			// If oneOffLabel only contains digits, the client typed something like:
			//   scalingo one-off-stop 1234
			labelHasOnlyDigit, err := regexp.MatchString("^[0-9]+$", oneOffLabel)
			if err != nil {
				// This should never occur as we are pretty sure the provided regexp is valid.
				errorQuit(c.Context, err)
			}
			if labelHasOnlyDigit {
				oneOffLabel = "one-off-" + oneOffLabel
			}

			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)

			err = apps.OneOffStop(c.Context, currentApp, oneOffLabel)
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "one-off-stop")
		},
	}
)
