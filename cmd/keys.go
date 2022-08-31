package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/keys"
)

var (
	listSSHKeyCommand = cli.Command{
		Name:     "keys",
		Category: "Public SSH Keys",
		Usage:    "List your SSH public keys",
		Description: `List all the public SSH keys associated with your account:

    $ scalingo keys

    # See also commands 'keys-add' and 'keys-remove'`,

		Action: func(c *cli.Context) error {
			err := keys.List(c.Context)
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "keys")
		},
	}

	addSSHKeyCommand = cli.Command{
		Name:     "keys-add",
		Category: "Public SSH Keys",
		Usage:    "Add a public SSH key to deploy your apps",
		Description: `Add a public SSH key:

    $ scalingo keys-add keyname /path/to/key

    # See also commands 'keys' and 'keys-remove'`,

		Action: func(c *cli.Context) error {
			if c.Args().Len() != 2 {
				cli.ShowCommandHelp(c, "keys-add")
				return nil
			}
			err := keys.Add(c.Context, c.Args().First(), c.Args().Slice()[1])
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "keys-add")
		},
	}

	removeSSHKeyCommand = cli.Command{
		Name:     "keys-remove",
		Category: "Public SSH Keys",
		Usage:    "Remove a public SSH key",
		Description: `Remove a public SSH key:

    $ scalingo keys-remove keyname

    # See also commands 'keys' and 'keys-add'`,

		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "keys-remove")
				return nil
			}
			err := keys.Remove(c.Context, c.Args().First())
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "keys-remove")
			autocomplete.KeysRemoveAutoComplete(c)
		},
	}
)
