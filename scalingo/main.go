package main

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/cmd"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/signals"
	"github.com/Scalingo/cli/update"
	"github.com/codegangsta/cli"
	"github.com/stvp/rollbar"
)

func main() {
	app := cli.NewApp()
	app.Name = "Scalingo Client"
	app.Author = "Scalingo Team"
	app.Email = "hello@scalingo.com"
	app.Usage = "Manage your apps and containers"
	app.Version = config.Version
	app.CategorizedHelp = true
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "app, a", Value: "<name>", Usage: "Name of the app", EnvVar: "SCALINGO_APP"},
	}
	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}
	app.Commands = []cli.Command{
		// Apps
		cmd.AppsCommand,
		cmd.CreateCommand,
		cmd.DestroyCommand,

		// Apps Actions
		cmd.LogsCommand,
		cmd.RunCommand,

		// Apps Process Actions
		cmd.PsCommand,
		cmd.ScaleCommand,
		cmd.RestartCommand,

		// Environment
		cmd.EnvCommand,
		cmd.EnvSetCommand,
		cmd.EnvUnsetCommand,

		// Domains
		cmd.DomainsListCommand,
		cmd.DomainsAddCommand,
		cmd.DomainsRemoveCommand,
		cmd.DomainsSSLCommand,

		// Addons
		cmd.AddonResourcesListCommand,
		cmd.AddonsListCommand,
		cmd.AddonsPlansCommand,

		// DB Access
		cmd.DbTunnelCommand,

		// SSH keys
		cmd.ListSSHKeyCommand,
		cmd.AddSSHKeyCommand,
		cmd.RemoveSSHKeyCommand,

		// Sessions
		cmd.LogoutCommand,
		cmd.SignUpCommand,

		// Version
		cmd.VersionCommand,
		cmd.UpdateCommand,
	}

	go signals.Handle()

	if len(os.Args) >= 2 && os.Args[1] == cmd.UpdateCommand.Name {
		err := update.Check()
		if err != nil {
			rollbar.Error(rollbar.ERR, err)
		}
		return
	} else {
		defer update.Check()
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("Fail to run scalingo", err)
	}
}
