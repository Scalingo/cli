package update

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Scalingo/cli/config"
	"gopkg.in/errgo.v1"
)

var (
	lastVersionURL = "https://raw.githubusercontent.com/Scalingo/appsdeck-executables/master/latest"
)

func Check() error {
	version := config.Version

	if strings.HasSuffix(version, "dev") {
		fmt.Println("No update checking, dev version:", version)
		return nil
	}

	lastVersion, err := getLastVersion()
	if err != nil {
		return errgo.Mask(err)
	}

	if version == lastVersion {
		return nil
	}

	fmt.Printf("Your Scalingo client (%s) is obsolete, some feature may not work correctly, please update to '%s'\n", version, lastVersion)
	return nil
}

func getLastVersion() (string, error) {
	res, err := http.Get(lastVersionURL)
	if err != nil {
		return "", errgo.Mask(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errgo.Mask(err)
	}

	return string(body), nil
}
