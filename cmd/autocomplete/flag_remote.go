package autocomplete

import (
	"context"
	"fmt"

	"github.com/Scalingo/cli/utils"
)

func FlagRemoteAutoComplete(ctx context.Context) bool {
	if dir, ok := utils.DetectGit(); ok {
		remoteNames := utils.ScalingoRepoAutoComplete(ctx, dir)
		for _, name := range remoteNames {
			fmt.Println(name)
		}
	} else {
		return false
	}
	return true
}
