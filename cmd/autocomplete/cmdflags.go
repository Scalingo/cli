package autocomplete

import (
	"os"

	"github.com/urfave/cli/v3"
)

func CmdFlagsAutoComplete(c *cli.Command, command string) error {
	if c.Name != command {
		return nil
	}

	if len(os.Args) > 1 {
		DisplayFlags(c.Flags)
	}

	return nil
}
