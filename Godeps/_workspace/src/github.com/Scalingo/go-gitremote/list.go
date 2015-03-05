package gitremote

import (
	"path/filepath"
	"regexp"
	"strings"
)

var (
	remoteTitleRe = regexp.MustCompile(`^\[remote\s"([A-Za-z0-9-_]+)"]$`)
	remoteURLRe   = regexp.MustCompile(`^\s+url\s*=\s*(.*)$`)
	remoteFetchRe = regexp.MustCompile(`^\s+fetch\s*=\s*(.*)$`)
	sshURLHostRe  = regexp.MustCompile(`^\w+@(.+):.+$`)
	sshURLRepoRe  = regexp.MustCompile(`^\w+@.+:(.+)$`)
)

// List returns the list of the remote for the GIT repository defined by path.
func List(path string) (Remotes, error) {
	configFilePath := filepath.Join(path, ".git", "config")
	config := &Config{Path: configFilePath}
	return config.List()
}

func ListFromConfigContent(content string) (Remotes, error) {
	var remotes Remotes
	var remote *Remote

	for _, line := range strings.Split(string(content), "\n") {
		matches := remoteTitleRe.FindStringSubmatch(line)
		if len(matches) == 2 {
			if remote != nil {
				remotes = append(remotes, remote)
				remote = nil
			}
			remote = &Remote{Name: matches[1]}
		}

		matches = remoteURLRe.FindStringSubmatch(line)
		if len(matches) == 2 {
			remote.URL = matches[1]
		}

		matches = remoteFetchRe.FindStringSubmatch(line)
		if len(matches) == 2 {
			remote.Fetch = matches[1]
		}
	}
	if remote != nil {
		remotes = append(remotes, remote)
	}

	return remotes, nil
}
