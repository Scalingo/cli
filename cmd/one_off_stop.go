package cmd

import (
	"context"
	"regexp"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-utils/errors/v2"
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
		Action: func(ctx context.Context, c *cli.Command) error {
			currentResource := detect.GetCurrentResource(ctx, c)
			if c.Args().Len() != 1 {
				_ = cli.ShowCommandHelp(ctx, c, "one-off-stop")
				return nil
			}
			isDB, err := utils.IsResourceDatabase(ctx, currentResource)
			if err != nil && !errors.Is(err, utils.ErrResourceNotFound) {
				errorQuit(ctx, err)
			}
			if isDB {
				io.Error("It is currently impossible to use `one-off-stop` on a database.")
				return nil
			}

			currentApp := currentResource

			oneOffLabel := c.Args().First()

			// If oneOffLabel only contains digits, the client typed something like:
			//   scalingo one-off-stop 1234
			labelHasOnlyDigit, err := regexp.MatchString("^[0-9]+$", oneOffLabel)
			if err != nil {
				// This should never occur as we are pretty sure the provided regexp is valid.
				errorQuit(ctx, err)
			}
			if labelHasOnlyDigit {
				oneOffLabel = "one-off-" + oneOffLabel
			}

			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)

			err = apps.OneOffStop(ctx, currentApp, oneOffLabel)
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "one-off-stop")
		},
	}
)
