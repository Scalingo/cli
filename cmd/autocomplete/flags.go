package autocomplete

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
)

func FlagsAutoComplete(c *cli.Context, flag string) bool {
	switch flag {
	case "-a", "--app":
		return FlagAppAutoComplete(c)
	}

	return false
}
