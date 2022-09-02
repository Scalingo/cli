package cmd

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/deployments"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/go-scalingo/v4/io"
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
				err := deployments.ResetCache(currentApp)
				if err != nil {
					errorQuit(err)
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
		Flags:    []cli.Flag{&appFlag},
		Description: ` List all of your previous app deployments
    $ scalingo -a myapp deployments
`,
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			err := deployments.List(currentApp)
			if err != nil {
				errorQuit(err)
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

			err := deployments.Logs(currentApp, deploymentID)
			if err != nil {
				errorQuit(err)
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
			err := deployments.Stream(&deployments.StreamOpts{
				AppName: currentApp,
			})
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
	}
	deploymentDeployCommand = cli.Command{
		Name:     "deploy",
		Category: "Deployment",
		Usage:    "Trigger a deployment by archive",
		Flags: []cli.Flag{&appFlag,
			&cli.BoolFlag{Name: "war", Aliases: []string{"w"}, Usage: "Specify that you want to deploy a WAR file"},
			&cli.BoolFlag{Name: "no-follow", Usage: "Return immediately after the deployment is triggered"},
		},
		Description: ` Trigger the deployment of a custom archive for your application

		scalingo deploy <archive path | archive URL> [version reference]

		The version reference is optional (generated from timestamp if none). It is a reference
		to the code you are deploying, version, commit SHA, etc.

		Examples:
		$ scalingo -a myapp deploy archive.tar.gz v1.0.0
		or
		$ scalingo -a myapp deploy http://example.com/archive.tar.gz v1.0.0
		or
		$ scalingo --app my-app deploy --no-follow archive.tar.gz v1.0.0
		$ scalingo --app my-app deployment-follow

    # See also commands 'deployments, deployment-follow'
`,
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
			opts := deployments.DeployOpts{NoFollow: c.Bool("no-follow")}
			if c.Bool("war") || strings.HasSuffix(archivePath, ".war") {
				io.Status(fmt.Sprintf("Deploying WAR archive: %s", archivePath))
				err := deployments.DeployWar(currentApp, archivePath, gitRef, opts)
				if err != nil {
					errorQuit(err)
				}
			} else {
				io.Status(fmt.Sprintf("Deploying tarball archive: %s", archivePath))
				err := deployments.Deploy(currentApp, archivePath, gitRef, opts)
				if err != nil {
					errorQuit(err)
				}
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "deploy")
		},
	}
)
