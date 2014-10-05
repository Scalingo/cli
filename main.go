package main

import (
	"fmt"
	"os"

	"github.com/Appsdeck/appsdeck/cmd"
	"github.com/Appsdeck/appsdeck/config"
	"github.com/Appsdeck/appsdeck/signals"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Appsdeck Client"
	app.Usage = "Manage your apps and containers"
	app.Version = config.Version
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "app, a", Value: "<name>", Usage: "Name of the app", EnvVar: "APPSDECK_APP"},
	}
	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}
	app.Commands = []cli.Command{
		cmd.ScaleCommand,
		cmd.LogsCommand, cmd.RunCommand,
		cmd.AppsCommand, cmd.LogoutCommand,
		cmd.CreateCommand, cmd.DestroyCommand,
		cmd.AddonsListCommand, cmd.AddonPlansCommand,
	}

	go signals.Handle()

	if err := app.Run(os.Args); err != nil {
		fmt.Println("Fail to run appsdeck", err)
	}
}
