package main

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/cmd"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/signals"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Scalingo Client"
	app.Usage = "Manage your apps and containers"
	app.Version = config.Version
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
		cmd.ScaleCommand,
		cmd.RestartCommand,

		// Environment
		cmd.EnvCommand,
		cmd.EnvSetCommand,
		cmd.EnvUnsetCommand,

		// Addons
		cmd.AddonsListCommand,
		cmd.AddonPlansCommand,
		cmd.AddonResourcesListCommand,

		// SSH keys
		cmd.ListSSHKeyCommand,
		cmd.AddSSHKeyCommand,
		cmd.RemoveSSHKeyCommand,

		// Misc
		cmd.LogoutCommand,
	}

	go signals.Handle()

	if err := app.Run(os.Args); err != nil {
		fmt.Println("Fail to run scalingo", err)
	}
}
