package autocomplete

import (
	"context"
	"encoding/gob"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v8"
	"github.com/Scalingo/go-scalingo/v8/debug"
)

type appsCache struct {
	CreatedAt time.Time
	Apps      []*scalingo.App
}

var (
	appsCacheFile     = filepath.Join(config.C.ConfigDir, ".apps-cache")
	appsCacheDuration = 30.0
	errExpiredCache   = errgo.New("apps has expired")
)

func appsList(ctx context.Context) ([]*scalingo.App, error) {
	var (
		err  error
		apps []*scalingo.App
	)

	apps, err = appsAutoCompleteCache()
	if err != nil {
		debug.Println("fail to get applications autocomplete cache make GET request", err)
		c, err := config.ScalingoClient(ctx)
		if err != nil {
			return nil, errgo.Notef(err, "fail to get Scalingo client")
		}
		apps, err = c.AppsList(ctx)
		if err != nil || len(apps) == 0 {
			return nil, errgo.Mask(err)
		}

		err = writeAppsAutoCompleteCache(apps)
		if err != nil {
			debug.Println("fail to write applications autocomplete cache", err)
		}
	}

	return apps, nil
}

func appsAutoCompleteCache() ([]*scalingo.App, error) {
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

	if time.Since(cache.CreatedAt).Seconds() > appsCacheDuration {
		return nil, errExpiredCache
	}

	return cache.Apps, nil
}

func writeAppsAutoCompleteCache(apps []*scalingo.App) error {
	fd, err := os.OpenFile(appsCacheFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return errgo.Mask(err)
	}
	defer func() { _ = fd.Close() }()
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
