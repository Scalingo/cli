package cmd

import (
	"os"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/db"
	"github.com/codegangsta/cli"
)

var (
	DbTunnelCommand = cli.Command{
		Name:     "db-tunnel",
		Category: "App Management",
		Usage:    "Create an encrypted connection to access your database",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "identity, i", Usage: "SSH Private Key", Value: os.Getenv("HOME") + "/.ssh/id_rsa", EnvVar: ""},
			cli.IntFlag{Name: "port, p", Usage: "Local port to bind", Value: 0, EnvVar: ""},
		},
		Description: `Create an SSH-encrypted connection to access your database locally
	Example
	  'scalingo --app my-app db-tunnel SCALINGO_MONGO_URL'`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "db-tunnel")
			} else if err := db.Tunnel(currentApp, c.Args()[0], c.String("identity"), c.Int("port")); err != nil {
				errorQuit(err)
			}
		},
	}
)
