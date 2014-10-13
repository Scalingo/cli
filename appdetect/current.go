package appdetect

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/debug"
)

func CurrentApp(appFlag string) (repoName string) {
	repoName = ""
	if appFlag != "<name>" {
		repoName = appFlag
	} else if DetectGit() {
		repoName, _ = ScalingoRepo()
	}
	if repoName == "" {
		fmt.Println("Unable to find the application name, please use --app flag.")
		os.Exit(1)
	} else {
		debug.Println("[AppDetect] App name is", repoName)
	}
	return
}
