package appdetect

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

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
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	file, err := os.OpenFile(path.Join(cwd, ".git", "config"), os.O_RDONLY, 0644)
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
			return found[2], nil
		}
	}
	return "", fmt.Errorf("Appsdeck GIT remote hasn't been found")
}
