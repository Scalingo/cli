package detect

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/utils"

	"github.com/Scalingo/go-scalingo/v4/debug"
)

// CurrentApp returns the app name if it has been found in one of the following:
// app flag, environment variable "SCALINGO_APP", current Git remote
// returns an empty string if not found
func CurrentApp(c *cli.Context) string {
	var appName string
	var err error

	if c.GlobalString("app") != "<name>" {
		appName = c.GlobalString("app")
	} else if c.String("app") != "<name>" {
		appName = c.String("app")
	} else if os.Getenv("SCALINGO_APP") != "" {
		appName = os.Getenv("SCALINGO_APP")
	} else if dir, ok := utils.DetectGit(); ok {
		appName, err = GetAppNameFromGitRemote(dir, c.GlobalString("remote"))
		if err != nil {
			debug.Println(err)
		}
	}
	if appName == "" {
		fmt.Println("Unable to find the application name, please use --app flag.")
		os.Exit(1)
	}
	debug.Println("[detect] App name is", appName)

	return appName
}

// GetAppNameFromGitRemote searches into the current directory and its parent for a remote
// named remoteName or scalingo-<remoteName>.
//
// It returns the application name and an error.
func GetAppNameFromGitRemote(directory string, remoteName string) (string, error) {
	remotes, err := utils.ScalingoGitRemotes(directory)
	if err != nil {
		return "", err
	}

	altRemoteName := "scalingo-" + remoteName
	for _, remote := range remotes {
		if remote.Config().Name == remoteName ||
			remote.Config().Name == altRemoteName {
			return extractAppNameFromGitRemote(remote.Config().URLs[0]), nil
		}
	}

	return "", errgo.Newf("[detect] Scalingo Git remote hasn't been found")
}

// extractAppNameFromGitRemote parses a Git remote and return the app name extracted
// out of it. The Git remote URL may look like:
// - SSH on a custom port: ssh://git@host:port/appName.git
// - GitHub: git@github.com:owner/appName.git
func extractAppNameFromGitRemote(url string) string {
	splittedURL := strings.Split(strings.TrimSuffix(url, ".git"), ":")
	appName := splittedURL[len(splittedURL)-1]
	// appName may contain "port/appName" or "owner/appName". We keep the part
	// after the slash.
	i := strings.LastIndex(appName, "/")
	if i != -1 {
		appName = appName[i+1:]
	}

	return appName
}
