package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/crontasks"
	"github.com/Scalingo/cli/detect"
)

var (
	cronTasksListCommand = cli.Command{
		Name:     "cron-tasks",
		Category: "Cron Tasks",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "List the cron tasks of an application",
		Description: CommandDescription{
			Description: "List all the cron tasks of an application",
			Examples:    []string{"scalingo --app my-app cron-tasks"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() > 0 {
				_ = cli.ShowCommandHelp(ctx, c, "cron-tasks")
				return nil
			}

			currentApp := detect.CurrentApp(c)
			err := crontasks.List(ctx, currentApp)
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "cron-tasks")
		},
	}
)
