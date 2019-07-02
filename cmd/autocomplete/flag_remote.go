package autocomplete

import (
	"fmt"

	"github.com/Scalingo/cli/appdetect"
	"github.com/urfave/cli"
)

func FlagRemoteAutoComplete(c *cli.Context) bool {
	if dir, ok := appdetect.DetectGit(); ok {
		remoteNames := appdetect.ScalingoRepoAutoComplete(dir)
		for _, name := range remoteNames {
			fmt.Println(name)
		}
	} else {
		return false
	}
	return true
}
