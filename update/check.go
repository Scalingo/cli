package update

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

var (
	lastVersionURL = "https://raw.githubusercontent.com/Scalingo/appsdeck-executables/master/latest"
	lastVersion    = ""
	gotLastVersion = make(chan struct{})
)

func init() {
	go func() {
		var err error
		lastVersion, err = getLastVersion()
		if err != nil {
			config.C.Logger.Println(err)
		}
		close(gotLastVersion)
	}()
}

func Check() error {
	version := config.Version

	if strings.HasSuffix(version, "dev") {
		fmt.Println("\nNo update checking, dev version:", version)
		return nil
	}

	<-gotLastVersion
	if version == lastVersion {
		return nil
	}

	io.Statusf(`Your Scalingo client (%s) is out-of-date: some features may not work correctly.
	Please update to '%s': https://github.com/Scalingo/cli/releases/tag/%s
`, version, lastVersion, lastVersion)
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

	return strings.TrimSpace(string(body)), nil
}
