package autocomplete

import (
	"os"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
)

func CmdFlagsAutoComplete(c *cli.Context, command string) error {

	var cmd cli.Command
	for _, cmd = range c.App.Commands {
		if cmd.Name == command {
			break
		}
	}
	if cmd.Name != command {
		return nil
	}

	if len(os.Args) > 1 {
		DisplayFlags(cmd.Flags)
	}

	return nil
}
