package autocomplete

import (
	"os"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/appdetect"
)

func CurrentAppCompletion(c *cli.Context) string {
	appName := ""
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
	if dir, ok := appdetect.DetectGit(); appName == "" && ok {
		appName, _ = appdetect.ScalingoRepo(dir, c.GlobalString("remote"))
	}
	return appName
}
