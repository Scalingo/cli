package main

import (
	"appsdeck/cli/cmd"
	"github.com/Appsdeck/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Appsdeck Client"
	app.Usage = "Manage your apps and containers"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{"app", "<name>", "Name of the app"},
	}
	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}
	app.Commands = []cli.Command{
		cmd.LogsCommand, cmd.RunCommand, cmd.AppsCommand, cmd.LogoutCommand,
	}

	app.Run(os.Args)
}
