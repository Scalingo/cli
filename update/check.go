package update

import (
	"encoding/json"
	"fmt"
	stdio "io"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/errgo.v1"

	"github.com/stvp/rollbar"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v4/debug"
)

const (
	lastVersionURL = "https://cli-dl.scalingo.com/version"
)

var (
	lastVersion    = ""
	cancelUpdate   = make(chan struct{})
	gotLastVersion = make(chan struct{})
	gotAnError     = false
)

type CLIVersionCache struct {
	Version string `json:"version"`
}

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

	if version != lastVersion {
		io.Errorf(io.BoldRed("Your Scalingo client (%s) is out-of-date: some features may not work correctly.\n"), version)
		io.Errorf(io.BoldRed("Please update to '%s' by reinstalling it: https://cli.scalingo.com\n"), lastVersion)
		return nil
	}

	return checkCLIVersionCache(version)
}

func checkCLIVersionCache(version string) error {
	// Check the cli version cache
	var cliVersionCache CLIVersionCache
	fd, err := os.Open(config.C.CLIVersionCachePath)
	if err == nil {
		err := json.NewDecoder(fd).Decode(&cliVersionCache)
		if err != nil {
			return errgo.Notef(err, "fail to decode cli version cache file")
		}
	}
	defer fd.Close()

	// This case happen if the cli has been upgraded
	if cliVersionCache.Version != "" && cliVersionCache.Version != version {
		// Show the changelog of each version since the last installed version.
		err := ShowChangelog(cliVersionCache.Version, version)
		if err != nil {
			rollbar.Error(rollbar.ERR, errgo.Notef(err, "fail to show last changelog"))
		}
	}

	// Save the version into a cache file.
	// To be able to track the update and show an accurate changelog.
	cliVersionCache.Version = version
	fd, err = os.OpenFile(config.C.CLIVersionCachePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0750)
	if err == nil {
		err := json.NewEncoder(fd).Encode(cliVersionCache)
		if err != nil {
			defer fd.Close()
			return errgo.Notef(err, "fail to encode cli version cache file")
		}
	} else {
		debug.Printf("[VERSION] Failed to save the cli version cache: %v\n", err)
	}
	defer fd.Close()

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
	body, err := stdio.ReadAll(res.Body)
	if err != nil {
		return "", errgo.Mask(err)
	}

	return strings.TrimSpace(string(body)), nil
}
