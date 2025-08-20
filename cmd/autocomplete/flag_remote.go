package autocomplete

import (
	"fmt"

	"github.com/Scalingo/cli/utils"
)

func FlagRemoteAutoComplete() bool {
	if dir, ok := utils.DetectGit(); ok {
		remoteNames := utils.ScalingoRepoAutoComplete(dir)
		for _, name := range remoteNames {
			fmt.Println(name)
		}
	} else {
		return false
	}
	return true
}
