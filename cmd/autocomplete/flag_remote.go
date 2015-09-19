package autocomplete

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/appdetect"
)

func FlagRemoteAutoComplete(c *cli.Context) bool {
	if dir, ok := appdetect.DetectGit(); ok {
		remoteNames := appdetect.ScalingoRepoComplete(dir)
		for _, name := range remoteNames {
			fmt.Println(name)
		}
	} else {
		return false
	}
	return true
}
