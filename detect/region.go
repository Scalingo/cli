package detect

import (
	"regexp"
	"strings"

	"github.com/urfave/cli/v3"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/utils"

	"github.com/Scalingo/go-scalingo/v8/debug"
)

// GetRegionFromGitRemote returns the region name extracted from remotes URL of the current Git repository
// If not found it returns an empty string
func GetRegionFromGitRemote(c *cli.Context, rc *config.RegionsCache) string {
	remoteName := RemoteNameFromFlags(c)

	if dir, ok := utils.DetectGit(); ok {
		remotes, err := utils.ScalingoGitRemotes(dir)
		if err != nil {
			debug.Println(err)
			return ""
		}

		altRemoteName := "scalingo-" + remoteName
		for _, remote := range remotes {
			if remote.Config().Name == remoteName || remote.Config().Name == altRemoteName {
				region, err := extractRegionFromGitRemote(remote.Config().URLs[0], rc)
				if err != nil {
					debug.Println(err)
				}
				debug.Println("[detect] region name is", region)
				return region
			}
		}
		debug.Println("[detect] fail to find a Scalingo Git remote")
		return ""
	}
	debug.Println("[detect] git repository cannot be found")
	return ""
}

// extractRegionFromGitRemote returns region name from a Git remote URL
// only if the region is a valid Scalingo region, returns an empty string othervise
func extractRegionFromGitRemote(url string, rc *config.RegionsCache) (string, error) {
	// Extract url part from the git remote
	// e.g.: git@ssh.osc-fr1.scalingo.com:app-name.git -> ssh.osc-fr1.scalingo.com
	r := regexp.MustCompile(`git@(.*)\:.*`)
	matches := r.FindStringSubmatch(url)
	if len(matches) > 0 {
		for _, region := range rc.Regions {
			if matches[1] == strings.Split(region.SSH, ":")[0] {
				return region.Name, nil
			}
		}
		return "", errgo.Newf("[detect] no valid region could be found in the Git remote")
	}
	return "", errgo.Newf("[detect] could not extract URL from Git remote")
}
