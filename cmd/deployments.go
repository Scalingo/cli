package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/deployments"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v8"
	"github.com/Scalingo/go-scalingo/v8/io"
)

var (
	deploymentCacheResetCommand = cli.Command{
		Name:     "deployment-delete-cache",
		Aliases:  []string{"deployment-cache-delete"},
		Category: "Deployment",
		Usage:    "Reset deployment cache",
		Flags:    []cli.Flag{&appFlag},
		Description: ` Delete the deployment cache (in case of corruption mostly)
    $ scalingo -a myapp deployment-delete-cache
`,
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() != 0 {
				_ = cli.ShowCommandHelp(ctx, c, "deployment-delete-cache")
			} else {
				currentApp := detect.CurrentApp(c)
				utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)
				err := deployments.ResetCache(ctx, currentApp)
				if err != nil {
					errorQuit(ctx, err)
				}
				io.Status("Deployment cache successfully deleted")
			}
			return nil
		},
	}

	deploymentsListCommand = cli.Command{
		Name:     "deployments",
		Category: "Deployment",
		Usage:    "List app deployments",
		Flags: []cli.Flag{
			&appFlag,
			&cli.IntFlag{Name: "page", Usage: "Page to display", Value: 1},
			&cli.IntFlag{Name: "per-page", Usage: "Number of deployments to display", Value: 20},
		},
		Description: ` List all of your previous app deployments
    $ scalingo -a myapp deployments
`,
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			err := deployments.List(ctx, currentApp, scalingo.PaginationOpts{
				Page:    c.Int("page"),
				PerPage: c.Int("per-page"),
			})
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
	}
	deploymentLogCommand = cli.Command{
		Name:     "deployment-logs",
		Category: "Deployment",
		Usage:    "View deployment logs",
		Flags:    []cli.Flag{&appFlag},
		Description: ` Get the logs of an app deployment
		$ scalingo -a myapp deployment-logs my-deployment
`,
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() > 1 {
				_ = cli.ShowCommandHelp(ctx, c, "deployment-logs")
			}

			deploymentID := ""
			if c.Args().Len() == 1 {
				deploymentID = c.Args().First()
			}

			err := deployments.Logs(ctx, currentApp, deploymentID)
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.DeploymentsAutoComplete(ctx, c)
		},
	}
	deploymentFollowCommand = cli.Command{
		Name:     "deployment-follow",
		Category: "Deployment",
		Usage:    "Follow deployment event stream",
		Flags:    []cli.Flag{&appFlag},
		Description: ` Get real-time deployment informations
		$ scalingo -a myapp deployment-follow
`,
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			err := deployments.Stream(ctx, &deployments.StreamOpts{
				AppName: currentApp,
			})
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
	}
	deploymentDeployCommand = cli.Command{
		Name:      "deploy",
		Category:  "Deployment",
		Usage:     "Trigger a deployment by archive",
		ArgsUsage: "<archive path | archive URL> [version reference]",
		Flags: []cli.Flag{&appFlag,
			&cli.BoolFlag{Name: "war", Aliases: []string{"w"}, Usage: "Specify that you want to deploy a WAR file"},
			&cli.BoolFlag{Name: "no-follow", Usage: "Return immediately after the deployment is triggered"},
		},
		Description: CommandDescription{
			Description: `Trigger the deployment of a custom archive for your application.

The version reference is optional (generated from timestamp if none).
It is a reference to the code you are deploying, version, commit SHA, etc.`,
			Examples: []string{
				"scalingo --app my-app deploy archive.tar.gz v1.0.0",
				"scalingo --app my-app deploy http://example.com/archive.tar.gz v1.0.0",
				"scalingo --app my-app deploy --no-follow archive.tar.gz v1.0.0",
				"scalingo --app my-app deployment-follow",
			},
			SeeAlso: []string{"deployments", "deployment-follow"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			args := c.Args()
			if args.Len() != 1 && args.Len() != 2 {
				_ = cli.ShowCommandHelp(ctx, c, "deploy")
				return nil
			}
			archivePath := args.First()
			gitRef := ""
			if args.Len() == 2 {
				gitRef = args.Slice()[1]
			}
			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)
			opts := deployments.DeployOpts{NoFollow: c.Bool("no-follow")}
			if c.Bool("war") || strings.HasSuffix(archivePath, ".war") {
				io.Status(fmt.Sprintf("Deploying WAR archive: %s", archivePath))
				err := deployments.DeployWar(ctx, currentApp, archivePath, gitRef, opts)
				if err != nil {
					errorQuit(ctx, err)
				}
			} else {
				io.Status(fmt.Sprintf("Deploying tarball archive: %s", archivePath))
				err := deployments.Deploy(ctx, currentApp, archivePath, gitRef, opts)
				if err != nil {
					errorQuit(ctx, err)
				}
			}
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "deploy")
		},
	}
)
