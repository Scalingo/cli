package cmd

import (
	"os"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/codegangsta-cli"
)

var (
	defaultKeyPath  = os.Getenv("HOME") + "/.ssh/id_rsa"
	DbTunnelCommand = cli.Command{
		Name:     "db-tunnel",
		Category: "App Management",
		Usage:    "Create an encrypted connection to access your database",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "identity, i", Usage: "SSH Private Key", Value: defaultKeyPath, EnvVar: ""},
			cli.IntFlag{Name: "port, p", Usage: "Local port to bind", Value: 0, EnvVar: ""},
		},
		Description: `Create an SSH-encrypted connection to access your database locally
	Example
	  'scalingo --app my-app db-tunnel SCALINGO_MONGO_URL'`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			var sshIdentity string
			if c.String("identity") == defaultKeyPath && os.Getenv("SSH_AUTH_SOCK") != "" {
				sshIdentity = "ssh-agent"
			} else {
				sshIdentity = c.String("identity")
			}
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "db-tunnel")
			} else if err := db.Tunnel(currentApp, c.Args()[0], sshIdentity, c.Int("port")); err != nil {
				errorQuit(err)
			}
		},
	}
)
