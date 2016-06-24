package cmd

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/keys"
	"github.com/Scalingo/codegangsta-cli"
	"github.com/mitchellh/go-homedir"
)

var (
	ListSSHKeyCommand = cli.Command{
		Name:     "keys",
		Category: "Public SSH Keys",
		Usage:    "List your SSH public keys",
		Description: `List all the public SSH keys associated with your account:

    $ scalingo keys

    # See also commands 'keys-add' and 'keys-remove'`,

		Action: func(c *cli.Context) {
			err := keys.List()
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "keys")
		},
	}

	AddSSHKeyCommand = cli.Command{
		Name:     "keys-add",
		Category: "Public SSH Keys",
		Usage:    "Add a public SSH key to deploy your apps",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "auto, a", Usage: "Autodetect key configuration"},
		},
		Description: `Add a public SSH key:

    $ scalingo keys-add keyname /path/to/key

		$ scalingo keys-add --auto

    # See also commands 'keys' and 'keys-remove'`,

		Action: func(c *cli.Context) {
			if len(c.Args()) != 2 && !c.Bool("auto") {
				cli.ShowCommandHelp(c, "keys-add")
				return
			}
			var keyname string
			var path string
			if c.Bool("auto") {
				var err error
				keyname, err = os.Hostname()
				if err != nil {
					errorQuit(err)
				}
				path, err = homedir.Expand("~/.ssh/id_rsa.pub")
				if err != nil {
					errorQuit(err)
				}

				fmt.Println("Adding key :")
				fmt.Println(" - Name: " + keyname)
				fmt.Println(" - Path: " + path)

			} else {
				keyname = c.Args()[0]
				path = c.Args()[1]
			}

			err := keys.Add(keyname, path)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "keys-add")
		},
	}

	RemoveSSHKeyCommand = cli.Command{
		Name:     "keys-remove",
		Category: "Public SSH Keys",
		Usage:    "Remove a public SSH key",
		Description: `Remove a public SSH key:

    $ scalingo keys-remove keyname

    # See also commands 'keys' and 'keys-add'`,

		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "keys-remove")
				return
			}
			err := keys.Remove(c.Args()[0])
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "keys-remove")
			autocomplete.KeysRemoveAutoComplete(c)
		},
	}
)
