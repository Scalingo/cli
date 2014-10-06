package appdetect

import (
	"github.com/Scalingo/cli/debug"
	"bufio"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

func openGitConfig() (*os.File, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path.Join(cwd, ".git", "config"), os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func DetectGit() bool {
	cwd, err := os.Getwd()
	if err != nil {
		return false
	}
	if _, err := os.Stat(path.Join(cwd, ".git")); err != nil {
		return false
	}
	return true
}

func AppsdeckRepo() (string, error) {
	file, err := openGitConfig()
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for line, err := reader.ReadString('\n'); err == nil; line, err = reader.ReadString('\n') {
		if strings.Contains(line, "url = git@appsdeck") {
			re := regexp.MustCompile(".*url = git@appsdeck.(eu|dev):(?P<repo>.*).git")
			found := re.FindStringSubmatch(line)
			if len(found) != 3 {
				return "", fmt.Errorf("Invalid Appsdeck GIT remote")
			}
			debug.Println("[AppDetect] GIT remote found:", strings.TrimSpace(line))
			return found[2], nil
		}
	}
	return "", fmt.Errorf("Appsdeck GIT remote hasn't been found")
}

func AddRemote(remote string) error {
	_, err := AppsdeckRepo()
	if err == nil {
		return fmt.Errorf("Remote already exists")
	}

	file, err := openGitConfig()
	if err != nil {
		return err
	}
	defer file.Close()
	file.Seek(0, os.SEEK_END)

	fmt.Fprintf(file,
		`[remote "appsdeck"]
	url = %s
	fetch = +refs/heads/*:refs/remotes/appsdeck/*
`, remote)

	return nil
}
