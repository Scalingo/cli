package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/deployments"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
)

var (
	DeploymentsListCommand = cli.Command{
		Name:     "deployments",
		Category: "Deployment",
		Usage:    "List app deployments",
		Flags:    []cli.Flag{appFlag},
		Description: ` List all of your previous app deployments
    $ scalingo -a myapp deployments
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			err := deployments.List(currentApp)
			if err != nil {
				errorQuit(err)
			}
		},
	}
	DeploymentLogCommand = cli.Command{
		Name:     "deployment-logs",
		Category: "Deployment",
		Usage:    "View deployment logs",
		Flags:    []cli.Flag{appFlag},
		Description: ` Get the logs of an app deployment
		$ scalingo -a myapp deployment-logs my-deployment
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) == 1 {
				err := deployments.Logs(currentApp, c.Args()[0])
				if err != nil {
					errorQuit(err)
				}
			} else {
				cli.ShowCommandHelp(c, "deployment-logs")
			}
		},
	}
	DeploymentFollowCommand = cli.Command{
		Name:     "deployment-follow",
		Category: "Deployment",
		Usage:    "Follow deployement event stream",
		Flags:    []cli.Flag{appFlag},
		Description: ` Get real-time deployment informations
		$ scalingo -a myapp deployment-follow
		`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			err := deployments.Stream(currentApp)
			if err != nil {
				errorQuit(err)
			}
		},
	}
)
