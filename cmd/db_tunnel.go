package cmd

import (
	"os"

	"github.com/urfave/cli"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/crypto/sshkeys"
	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/cli/detect"
)

var (
	DbTunnelCommand = cli.Command{
		Name:     "db-tunnel",
		Category: "App Management",
		Usage:    "Create an encrypted connection to access your database",
		Flags: []cli.Flag{appFlag,
			cli.IntFlag{Name: "port, p", Usage: "Local port to bind (default 10000)"},
			cli.StringFlag{Name: "identity, i", Usage: "SSH Private Key"},
			cli.StringFlag{Name: "bind, b", Usage: "IP to bind (default 127.0.0.1)"},
			cli.BoolTFlag{Name: "reconnect", Usage: "true by default, automatically reconnect to the tunnel when disconnected"},
		},
		Description: `Create an SSH-encrypted connection to access your Scalingo database locally.

   This command works for all the database addons provided by Scalingo. MySQL,
   PostgreSQL, MongoDB, Redis or Elasticsearch. This action authenticates you
   thanks to your SSH key (exactly the same as a 'git push' operation).

   The command takes one argument which is, either the name of an environment
   variable of your app, or its value containing the connection URL to your
   database.

   Example
     $ scalingo --app my-app db-tunnel SCALINGO_MONGO_URL
     $ scalingo --app my-app db-tunnel mongodb://user:pass@host:port/db

   Once the tunnel is built, the port which has been allocated will be
   displayed (default is 10000), example: "localhost:10000". You can
	 choose this port manually with the '-p' option.

   Example
     $ scalingo --app my-app db-tunnel -p 20000 MONGO_URL
     Building tunnel to my-app-1.mongo.dbs.scalingo.com:12345
     You can access your database on '127.0.0.1:20000'

     $ mongo -u <user> -p <pass> 127.0.0.1:20000/my-app-1
     >

   We are looking if an SSH-agent is running on your host, otherwise we are
   using the SSH key '$HOME/.ssh/id_rsa'. You can specify a precise SSH key
   you want to use to authenticate thanks to the '-i' flag.

   Example
     $ scalingo --app rails-app db-tunnel -i ~/.ssh/custom_key DATABASE_URL`,
		Action: func(c *cli.Context) {
			currentApp := detect.CurrentApp(c)
			var sshIdentity string
			if c.String("identity") == "" && os.Getenv("SSH_AUTH_SOCK") != "" {
				sshIdentity = "ssh-agent"
			} else if c.String("identity") == "" {
				sshIdentity = sshkeys.DefaultKeyPath
			} else {
				sshIdentity = c.String("identity")
			}
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "db-tunnel")
				return
			}

			err := db.Tunnel(db.TunnelOpts{
				App:       currentApp,
				DBEnvVar:  c.Args()[0],
				Identity:  sshIdentity,
				Port:      c.Int("port"),
				Bind:      c.String("bind"),
				Reconnect: c.BoolT("reconnect"),
			})
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "db-tunnel")
			autocomplete.DbTunnelAutoComplete(c)
		},
	}
)
