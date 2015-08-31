package autocomplete

import (
	"fmt"
	"os"
	//	"reflect"
	"strings"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
)

func RunAutoComplete(c *cli.Context) error {

	var cmd cli.Command
	for _, cmd = range c.App.Commands {
		if cmd.Name == "run" {
			break
		}
	}
	if cmd.Name != "run" {
		return nil
	}

	if len(os.Args) > 1 {
		lastArg := os.Args[len(os.Args)-2]

		found := false
		var lastFlag cli.Flag
		for _, lastFlag = range cmd.Flags {
			names := GetFlagNames(lastFlag)
			for i := range names {
				if names[i] == lastArg {
					found = true
					break
				}
			}
		}
		fmt.Printf("%s===%s===%T\n", lastArg, lastFlag, lastFlag)
		if !strings.HasPrefix(lastArg, "-") || found {
			DisplayFlags(cmd.Flags)
		}
	}
	return nil
}
