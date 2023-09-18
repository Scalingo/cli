package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config/auth"
	"github.com/Scalingo/go-scalingo/v6"
	"github.com/Scalingo/go-scalingo/v6/debug"
)

type UnknownRegionError struct {
	regionName string
}

func (err UnknownRegionError) Error() string {
	return fmt.Sprintf("invalid region %v", err.regionName)
}

type RegionsCache struct {
	ExpireAt time.Time         `json:"expire_at"`
	Regions  []scalingo.Region `json:"regions"`
}

func (c RegionsCache) Default() (scalingo.Region, error) {
	if len(c.Regions) == 0 {
		return scalingo.Region{}, fmt.Errorf("no region found")
	}

	for _, r := range c.Regions {
		if r.Default {
			return r, nil
		}
	}

	return scalingo.Region{}, fmt.Errorf("no default region found")
}

// GetRegionOpts allows the caller to use a custom API token instead of
// the default one of the authentication used
type GetRegionOpts struct {
	Token    string
	SkipAuth bool
}

func EnsureRegionsCache(ctx context.Context, c Config, opts GetRegionOpts) (RegionsCache, error) {
	debug.Println("[Regions] Ensure cache is filled")
	var regionsCache RegionsCache
	fd, err := os.Open(c.RegionsCachePath)
	if err != nil {
		debug.Printf("[Regions] Fail to open the cache: %v\n", err)
	} else {
		json.NewDecoder(fd).Decode(&regionsCache)
		fd.Close()
	}
	// If cache is still valid
	if time.Now().Unix() <= regionsCache.ExpireAt.Unix() {
		debug.Println("[Regions] Use the cache")
		return regionsCache, nil
	}

	var client *scalingo.Client
	if opts.SkipAuth {
		debug.Println("[Regions] Create an unauthenticated client to the authentication service")
		client, err = ScalingoUnauthenticatedAuthClient(ctx)
		if err != nil {
			return RegionsCache{}, errgo.Notef(err, "fail to create an unauthenticated client")
		}
	} else {
		token := &auth.UserToken{Token: opts.Token}
		if token.Token == "" {
			auth := &CliAuthenticator{}
			_, token, err = auth.LoadAuth()
			if err != nil {
				return RegionsCache{}, errgo.Notef(err, "fail to load authentication")
			}
		}

		debug.Println("[Regions] Create an authenticated client to the authentication service using the token")
		client, err = ScalingoAuthClientFromToken(ctx, token.Token)
		if err != nil {
			return RegionsCache{}, errgo.Notef(err, "fail to create an authenticated Scalingo client using the API token")
		}
	}

	debug.Println("[Regions] Get the list of regions to fill the cache")
	regions, err := client.RegionsList(ctx)
	if err != nil {
		return RegionsCache{}, errgo.Notef(err, "fail to list available regions")
	}

	regionsCache.Regions = regions
	regionsCache.ExpireAt = time.Now().Add(10 * time.Minute)

	if opts.SkipAuth {
		debug.Println("[Regions] Do not save the cache")
		// If we skipped the authentication the region cache should not be saved since it will not contain regions that are not publicly available (like osc-secnum-fr1)
		return regionsCache, nil
	}

	debug.Println("[Regions] Save the cache")
	fd, err = os.OpenFile(c.RegionsCachePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0750)
	if err != nil {
		debug.Printf("[Regions] Failed to save the cache: %v\n", err)
		return regionsCache, nil
	}

	json.NewEncoder(fd).Encode(regionsCache)
	fd.Close()
	return regionsCache, nil
}

// GetRegion returns the requested region configuration, use local file system
// cache if any. In case of cache fault, save on disk for 10 minutes the
// available regions
func GetRegion(ctx context.Context, c Config, name string, opts GetRegionOpts) (scalingo.Region, error) {
	regionsCache, err := EnsureRegionsCache(ctx, c, opts)
	if err != nil {
		return scalingo.Region{}, errgo.Notef(err, "fail to get the regions cache")
	}

	for _, region := range regionsCache.Regions {
		if region.Name == name {
			debug.Println("[Regions] Found the region in the cache")
			return region, nil
		}
	}
	return scalingo.Region{}, UnknownRegionError{regionName: name}
}
