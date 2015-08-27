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
	remotes, err := gitremote.List(directory)
	if err != nil {
		return "", err
	}
	for i := 0; i < 2; i++ {
		for _, remote := range remotes {
			if remote.Name == remoteName {
				matched, err := regexp.Match(".*scalingo.com:.*.git", []byte(remote.URL))
				if err == nil && matched {
					debug.Println("[AppDetect] GIT remote found:", remote)
					return filepath.Base(strings.TrimSuffix(remote.Repository(), ".git")), nil
				}
			}
		}
		remoteName = "scalingo-" + remoteName
	}
	return "", errgo.Newf("Scalingo GIT remote hasn't been found")
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
