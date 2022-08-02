package utils

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

// DetectGit detects if current directory is a Git repository
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

func ScalingoRepoAutoComplete(dir string) []string {
	var repos []string

	remotes, err := ScalingoGitRemotes(dir)
	if err != nil {
		debug.Println("[detectCompletion] fail to get scalingo remotes in", dir)
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

// ScalingoGitRemotes returns an array of remote URLs (*.scalingo.com) from a git repository <directory>
func ScalingoGitRemotes(directory string) ([]*git.Remote, error) {
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
		if len(remote.Config().URLs) == 0 {
			continue
		}

		remoteURL := remote.Config().URLs[0]
		matched, err := regexp.Match(".*scalingo.com:.*.git", []byte(remoteURL))
		if err != nil || !matched {
			continue
		}

		debug.Println("git remote found:", remoteURL)
		matchedRemotes = append(matchedRemotes, remote)
	}

	return matchedRemotes, nil
}

// AddGitRemote add a remote URL to the current git repository
func AddGitRemote(url string, name string) error {
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
