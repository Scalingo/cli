package autocomplete

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/debug"
	"gopkg.in/errgo.v1"
)

type appsCache struct {
	CreatedAt time.Time
	Apps      []*api.App
}

var (
	appsCacheFile     = filepath.Join(config.C.ConfigDir, ".apps-cache")
	appsCacheDuration = 10.0
	errExpiredCache   = errgo.New("apps has expired")
)

func FlagAppAutoComplete(c *cli.Context) bool {
	var (
		err  error
		apps []*api.App
	)

	apps, err = appsAutoCompleteCache()
	if err != nil {
		debug.Println("fail to get applications autocomplete cache", err)
		apps, err = api.AppsList()
		if err != nil || len(apps) == 0 {
			return false
		}

		err = writeAppsAutoCompleteCache(apps)
		if err != nil {
			debug.Println("fail to write applications autocomplete cache", err)
		}
	}

	for _, app := range apps {
		fmt.Println(app.Name)
	}

	return true
}

func appsAutoCompleteCache() ([]*api.App, error) {
	fd, err := os.Open(appsCacheFile)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer fd.Close()
	var cache appsCache
	err = gob.NewDecoder(fd).Decode(&cache)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	if time.Now().Sub(cache.CreatedAt).Seconds() > appsCacheDuration {
		return nil, errExpiredCache
	}

	return cache.Apps, nil
}

func writeAppsAutoCompleteCache(apps []*api.App) error {
	fd, err := os.OpenFile(appsCacheFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return errgo.Mask(err)
	}
	defer fd.Close()
	cache := appsCache{
		Apps:      apps,
		CreatedAt: time.Now(),
	}
	err = gob.NewEncoder(fd).Encode(&cache)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}
