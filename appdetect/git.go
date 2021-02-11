package appdetect

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v4/debug"
)

func DetectGit() (string, bool) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", false
	}
	for cwd != "/" {
		if _, err := os.Stat(path.Join(cwd, ".git")); err == nil {
			return cwd, true
		}
		cwd = filepath.Dir(cwd)
	}
	return "", false
}

// ScalingoRepo searches into the current directory and its parent for a remote
// named remoteName or scalingo-<remoteName>.
//
// It returns the application name and an error.
func ScalingoRepo(directory string, remoteName string) (string, error) {
	remotes, err := scalingoRemotes(directory)
	if err != nil {
		return "", err
	}

	altRemoteName := "scalingo-" + remoteName
	for _, remote := range remotes {
		if remote.Config().Name == remoteName ||
			remote.Config().Name == altRemoteName {
			return getAppNameFromGitRemote(remote.Config().URLs[0]), nil
		}
	}
	return "", errgo.Newf("Scalingo Git remote hasn't been found")
}

// getAppNameFromGitRemote parses a Git remote and return the app name extracted
// out of it. The Git remote URL may look like:
// - agora-fr1: git@host:appName.git
// - SSH on a custom port: ssh://git@host:port/appName.git
// - GitHub: git@github.com:owner/appName.git
func getAppNameFromGitRemote(url string) string {
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

func ScalingoRepoAutoComplete(dir string) []string {
	var repos []string

	remotes, err := scalingoRemotes(dir)
	if err != nil {
		debug.Println("[AppDetectCompletion] fail to get scalingo remotes in", dir)
		return repos
	}

	for _, remote := range remotes {
		if strings.HasPrefix(remote.Config().Name, "scalingo-") {
			repos = append(repos, remote.Config().Name[9:])
		} else {
			repos = append(repos, remote.Config().Name)
		}
	}

	return repos
}

func scalingoRemotes(directory string) ([]*git.Remote, error) {
	repo, err := git.PlainOpen(directory)
	if err != nil {
		return nil, errgo.Notef(err, "fail to initialize the Git repository")
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return nil, errgo.Notef(err, "fail to list the remotes")
	}

	matchedRemotes := []*git.Remote{}
	for _, remote := range remotes {
		remoteURL := remote.Config().URLs[0]
		matched, err := regexp.Match(".*scalingo.com:.*.git", []byte(remoteURL))
		if err != nil || !matched {
			continue
		}

		debug.Println("[AppDetect] Git remote found:", remoteURL)
		matchedRemotes = append(matchedRemotes, remote)
	}

	return matchedRemotes, nil
}

func AddRemote(url string, name string) error {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return errgo.Notef(err, "fail to initialize the Git repository")
	}

	_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
		Name: name,
		URLs: []string{url},
	})
	if err != nil {
		return errgo.Notef(err, "fail to add the Git remote")
	}

	return nil
}
