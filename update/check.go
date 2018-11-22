package update

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

var (
	lastVersionURL = "https://cli-dl.scalingo.com/version"
	lastVersion    = ""
	cancelUpdate   = make(chan struct{})
	gotLastVersion = make(chan struct{})
	gotAnError     = false
)

func init() {
	if config.C.DisableUpdateChecker {
		close(cancelUpdate)
		return
	}
	go func() {
		var err error
		lastVersion, err = getLastVersion()
		if err != nil {
			config.C.Logger.Println(err)
			gotAnError = true
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

	select {
	case <-cancelUpdate:
		return nil
	case <-gotLastVersion:
	}

	if gotAnError {
		return errgo.New("Update checker: connection error")
	}
	if version == lastVersion {
		return nil
	}

	io.Statusf(io.BoldRed("Your Scalingo client (%s) is out-of-date: some features may not work correctly.\n"), version)
	io.Infof(io.BoldRed("Please update to '%s' by reinstalling it: https://cli.scalingo.com\n"), lastVersion)
	return nil
}

func getLastVersion() (string, error) {
	client := http.Client{
		Timeout: 4 * time.Second,
	}

	res, err := client.Get(lastVersionURL)
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
