package autocomplete

import (
	"context"
	"encoding/gob"
	stderrors "errors"
	"os"
	"path/filepath"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-scalingo/v11/debug"
	"github.com/Scalingo/go-utils/errors/v3"
)

type appsCache struct {
	CreatedAt time.Time
	Apps      []*scalingo.App
}

var (
	appsCacheFile     = filepath.Join(config.C.ConfigDir, ".apps-cache")
	appsCacheDuration = 30.0
	errExpiredCache   = stderrors.New("apps has expired")
)

func appsList(ctx context.Context) ([]*scalingo.App, error) {
	var (
		err  error
		apps []*scalingo.App
	)

	apps, err = appsAutoCompleteCache(ctx)
	if err != nil {
		debug.Println("fail to get applications autocomplete cache make GET request", err)
		c, err := config.ScalingoClient(ctx)
		if err != nil {
			return nil, errors.Wrapf(ctx, err, "fail to get Scalingo client")
		}
		apps, err = c.AppsList(ctx)
		if err != nil || len(apps) == 0 {
			return nil, errors.Wrap(ctx, err, "list applications")
		}

		err = writeAppsAutoCompleteCache(ctx, apps)
		if err != nil {
			debug.Println("fail to write applications autocomplete cache", err)
		}
	}

	return apps, nil
}

func appsAutoCompleteCache(ctx context.Context) ([]*scalingo.App, error) {
	fd, err := os.Open(appsCacheFile)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "open apps cache file %s", appsCacheFile)
	}
	defer fd.Close()
	var cache appsCache
	err = gob.NewDecoder(fd).Decode(&cache)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "decode apps cache file")
	}

	if time.Since(cache.CreatedAt).Seconds() > appsCacheDuration {
		return nil, errExpiredCache
	}

	return cache.Apps, nil
}

func writeAppsAutoCompleteCache(ctx context.Context, apps []*scalingo.App) error {
	fd, err := os.OpenFile(appsCacheFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return errors.Wrapf(ctx, err, "open apps cache file %s for write", appsCacheFile)
	}
	defer func() { _ = fd.Close() }()
	cache := appsCache{
		Apps:      apps,
		CreatedAt: time.Now(),
	}
	err = gob.NewEncoder(fd).Encode(&cache)
	if err != nil {
		return errors.Wrap(ctx, err, "encode apps cache file")
	}

	return nil
}
