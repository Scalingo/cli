package appdetect

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/go-gitremote"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/debug"
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

func ScalingoRepo(directory string, remoteName string) (string, error) {
	remotes, err := scalingoRemotes(directory)
	if err != nil {
		return "", err
	}

	altRemoteName := "scalingo-" + remoteName
	for _, remote := range remotes {
		if remote.Name == remoteName || remote.Name == altRemoteName {
			return filepath.Base(strings.TrimSuffix(remote.Repository(), ".git")), nil
		}
	}
	return "", errgo.Newf("Scalingo GIT remote hasn't been found")
}

func ScalingoRepoComplete(dir string) []string {
	var repos []string

	remotes, err := scalingoRemotes(dir)
	if err != nil {
		debug.Println("[AppDetectCompletion] fail to get scalingo remotes in", dir)
		return repos
	}

	for _, remote := range remotes {
		if strings.HasPrefix(remote.Name, "scalingo-") {
			repos = append(repos, remote.Name[9:])
		} else {
			repos = append(repos, remote.Name)
		}
	}

	return repos
}

func scalingoRemotes(directory string) (gitremote.Remotes, error) {
	matchedRemotes := make(gitremote.Remotes, 0)
	remotes, err := gitremote.List(directory)
	if err != nil {
		return nil, err
	}
	for _, remote := range remotes {
		matched, err := regexp.Match(".*scalingo.com:.*.git", []byte(remote.URL))
		if err == nil && matched {
			debug.Println("[AppDetect] GIT remote found:", remote)
			matchedRemotes = append(matchedRemotes, remote)
		}
	}
	return matchedRemotes, nil
}

func AddRemote(url string, name string) error {
	remote := &gitremote.Remote{
		Name: name,
		URL:  url,
	}

	config := gitremote.New(".")
	err := config.Add(remote)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}
