package appdetect

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/Scalingo/cli/debug"
)

func CurrentApp(c *cli.Context) string {
	var repoName string
	if c.GlobalString("app") != "<name>" {
		repoName = c.GlobalString("app")
	} else if c.String("app") != "<name>" {
		repoName = c.String("app")
	} else if os.Getenv("SCALINGO_APP") != "" {
		repoName = os.Getenv("SCALINGO_APP")
	} else if dir, ok := DetectGit(); ok {
		repoName, _ = ScalingoRepo(dir, c.GlobalString("remote"))
	}
	if repoName == "" {
		fmt.Println("Unable to find the application name, please use --app flag.")
		os.Exit(1)
	}

	debug.Println("[AppDetect] App name is", repoName)
	return repoName
}
