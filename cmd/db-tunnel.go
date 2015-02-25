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
		Flags: []cli.Flag{appFlag,
			cli.StringFlag{Name: "identity, i", Usage: "SSH Private Key", Value: defaultKeyPath, EnvVar: ""},
			cli.IntFlag{Name: "port, p", Usage: "Local port to bind", Value: 0, EnvVar: ""},
		},
		Description: `Create an SSH-encrypted connection to access your database locally. This
   action authenticate you thanks to your SSH key (exactly the same as a 'git
   push' operation).

   We are looking if an SSH-agent is running on your host, otherwise we are
   using the SSH key '$HOME/.ssh/id_rsa'. You can specify a precise SSH key
   you want to use to authenticate thanks to the '-i' flag.

   The command take one argument which is, either the name of an environment
   variable or it value, containing the connection URL to your database.

   Then, a tunnel is built and the command will display a port which has been
   allocated for local usage, example: "localhost:58000", if you want the
   database with a specific port, user the '-p' option.

   Example
     'scalingo -a my-app db-tunnel SCALINGO_MONGO_URL'
     'scalingo -a rails-app db-tunnel -i ~/.ssh/custom_key -p 5432 DATABASE_URL'`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
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
