package appdetect

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/Scalingo/cli/debug"
	"gopkg.in/errgo.v1"
)

func openGitConfig() (*os.File, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	file, err := os.OpenFile(path.Join(cwd, ".git", "config"), os.O_RDWR, 0644)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
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

func ScalingoRepo() (string, error) {
	file, err := openGitConfig()
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for line, err := reader.ReadString('\n'); err == nil; line, err = reader.ReadString('\n') {
		if strings.Contains(line, "url = git@appsdeck") || strings.Contains(line, "url = git@scalingo") {
			re := regexp.MustCompile(".*url = git@(appsdeck|scalingo).(com|eu|dev):(?P<repo>.*).git")
			found := re.FindStringSubmatch(line)
			if len(found) != 3 {
				return "", errgo.Newf("Invalid ScalingoGIT remote")
			}
			debug.Println("[AppDetect] GIT remote found:", strings.TrimSpace(line))
			return found[2], nil
		}
	}
	return "", errgo.Newf("Scalingo GIT remote hasn't been found")
}

func AddRemote(remote string) error {
	_, err := ScalingoRepo()
	if err == nil {
		return errgo.Notef(err, "remote already exists")
	}

	file, err := openGitConfig()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer file.Close()
	file.Seek(0, os.SEEK_END)

	fmt.Fprintf(file,
		`[remote "scalingo"]
	url = %s
	fetch = +refs/heads/*:refs/remotes/scalingo/*
`, remote)

	return nil
}
