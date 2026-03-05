package update

import (
	"context"
	"fmt"
	stdio "io"
	"net/http"
	"strings"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-utils/errors/v3"
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

func init() {
	if config.C.DisableUpdateChecker {
		close(cancelUpdate)
		return
	}
	go func() {
		var err error
		lastVersion, err = getLastVersion(context.Background())
		if err != nil {
			config.C.Logger.Println(err)
			gotAnError = true
		}
		close(gotLastVersion)
	}()
}

func Check(ctx context.Context) error {
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
		return errors.New(ctx, "Update checker: connection error")
	}
	if version == lastVersion {
		return nil
	}

	io.Errorf(io.BoldRed("Your Scalingo client (%s) is out-of-date: some features may not work correctly.\n"), version)
	io.Errorf(io.BoldRed("Please update to '%s' by reinstalling it: https://cli.scalingo.com\n"), lastVersion)
	return nil
}

func getLastVersion(ctx context.Context) (string, error) {
	client := http.Client{
		Timeout: 4 * time.Second,
	}

	res, err := client.Get(lastVersionURL)
	if err != nil {
		return "", errors.Wrap(ctx, err, "fetch latest CLI version")
	}
	defer res.Body.Close()
	body, err := stdio.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(ctx, err, "read latest CLI version response")
	}

	return strings.TrimSpace(string(body)), nil
}
