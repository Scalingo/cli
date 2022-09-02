package autocomplete

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v5/debug"
)

func CurrentAppCompletion(c *cli.Context) string {
	appName := ""
	var err error

	if len(os.Args) >= 2 {
		for i := range os.Args {
			if i < len(os.Args) && (os.Args[i] == "-a" || os.Args[i] == "--app") {
				if (i+1) < len(os.Args) && os.Args[i+1] != "" {
					return os.Args[i+1]
				}
			}
		}
	}
	if appName == "" && os.Getenv("SCALINGO_APP") != "" {
		appName = os.Getenv("SCALINGO_APP")
	}
	if dir, ok := utils.DetectGit(); appName == "" && ok {
		appName, err = detect.GetAppNameFromGitRemote(dir, detect.RemoteNameFromFlags(c))
		if err != nil {
			debug.Println(err)
		}
	}
	return appName
}
