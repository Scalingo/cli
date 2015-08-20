package autocomplete

import (
	"os"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/appdetect"
)

func CurrentAppCompletion(c *cli.Context) string {
	appName := ""
	if len(os.Args) >= 2 {
		for a := range os.Args {
			if a < len(os.Args) && (os.Args[a] == "-a" || os.Args[a] == "-app") {
				if (a + 1) < len(os.Args) {
					appName = os.Args[a+1]
				}
			}
		}
	} else if appdetect.DetectGit() {
		appName, _ = appdetect.ScalingoRepo()
	} else if os.Getenv("SCALINGO_APP") != "" {
		appName = os.Getenv("SCALINGO_APP")
	}
	return appName
}
