package appdetect

import (
	"fmt"
	"os"
)

func CurrentApp(appFlag string) string {
	if appFlag != "<name>" {
		return appFlag
	} else if DetectGit() {
		repoName, err := AppsdeckRepo()
		if err == nil {
			return repoName
		}
	}
	fmt.Println("Unable to find the repository name, please use --app flag.")
	os.Exit(1)
	return ""
}
