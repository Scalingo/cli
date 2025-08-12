package cmd

import (
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
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "deployment-delete-cache")
			} else {
				currentApp := detect.CurrentApp(c)
				utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)
				err := deployments.ResetCache(c.Context, currentApp)
				if err != nil {
					errorQuit(c.Context, err)
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
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			err := deployments.List(c.Context, currentApp, scalingo.PaginationOpts{
				Page:    c.Int("page"),
				PerPage: c.Int("per-page"),
			})
			if err != nil {
				errorQuit(c.Context, err)
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
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() > 1 {
				cli.ShowCommandHelp(c, "deployment-logs")
			}

			deploymentID := ""
			if c.Args().Len() == 1 {
				deploymentID = c.Args().First()
			}

			err := deployments.Logs(c.Context, currentApp, deploymentID)
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.DeploymentsAutoComplete(c)
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
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			err := deployments.Stream(c.Context, &deployments.StreamOpts{
				AppName: currentApp,
			})
			if err != nil {
				errorQuit(c.Context, err)
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

		Action: func(c *cli.Context) error {
			args := c.Args()
			if args.Len() != 1 && args.Len() != 2 {
				cli.ShowCommandHelp(c, "deploy")
				return nil
			}
			archivePath := args.First()
			gitRef := ""
			if args.Len() == 2 {
				gitRef = args.Slice()[1]
			}
			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)
			opts := deployments.DeployOpts{NoFollow: c.Bool("no-follow")}
			if c.Bool("war") || strings.HasSuffix(archivePath, ".war") {
				io.Status(fmt.Sprintf("Deploying WAR archive: %s", archivePath))
				err := deployments.DeployWar(c.Context, currentApp, archivePath, gitRef, opts)
				if err != nil {
					errorQuit(c.Context, err)
				}
			} else {
				io.Status(fmt.Sprintf("Deploying tarball archive: %s", archivePath))
				err := deployments.Deploy(c.Context, currentApp, archivePath, gitRef, opts)
				if err != nil {
					errorQuit(c.Context, err)
				}
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "deploy")
		},
	}
)
