package cmd

import (
	"context"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/go-utils/errors/v2"
)

var (
	appsCommand = cli.Command{
		Name:        "apps",
		Category:    "Global",
		Description: "List your apps and give some details about them",
		Flags:       []cli.Flag{&cli.StringFlag{Name: "project", Usage: "Filter apps by project. The filter uses the format <ownerUsername>/<projectName>"}},
		Usage:       "List your apps",
		Action: func(ctx context.Context, c *cli.Command) error {
			projectSlug := c.String("project")
			if projectSlug != "" {
				projectSlugSplit := strings.Split(projectSlug, "/")
				if len(projectSlugSplit) != 2 || (len(projectSlugSplit) == 2 && (projectSlugSplit[0] == "" || projectSlugSplit[1] == "")) {
					errorQuitWithHelpMessage(errors.New(ctx, "project filter doesn't respect the expected format"), c, "apps")
				}
			}
			if err := apps.List(ctx, projectSlug); err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "apps")
		},
	}

	appsInfoCommand = cli.Command{
		Name:     "apps-info",
		Category: "App Management",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "Display the application information",
		Description: CommandDescription{
			Description: "Display various application information such as the force HTTPS status, the stack configured, sticky sessions, etc.",
			Examples:    []string{"scalingo apps-info --app my-app"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if err := apps.Info(ctx, currentApp); err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "apps-info")
		},
	}
)
