package autocomplete

import (
	"os"

	"github.com/urfave/cli/v3"
)

func CmdFlagsAutoComplete(c *cli.Command, command string) error {
	var cmd *cli.Command
	for _, cmd = range c.Commands {
		if cmd.Name == command {
			break
		}
	}
	if cmd == nil || cmd.Name != command {
		return nil
	}

	if len(os.Args) > 1 {
		DisplayFlags(cmd.Flags)
	}

	return nil
}
