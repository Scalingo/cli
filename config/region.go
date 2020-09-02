package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Scalingo/cli/config/auth"
	"github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/go-scalingo/debug"
	"gopkg.in/errgo.v1"
)

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

func EnsureRegionsCache(c Config, opts GetRegionOpts) (RegionsCache, error) {
	var regionsCache RegionsCache
	fd, err := os.Open(c.RegionsCachePath)
	if err == nil {
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
		client, err = ScalingoUnauthenticatedAuthClient()
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

		debug.Println("[Regions] Get the list of regions to fill the cache")
		client, err = ScalingoAuthClientFromToken(token.Token)
		if err != nil {
			return RegionsCache{}, errgo.Notef(err, "fail to create an authenticated Scalingo client using the API token")
		}
	}

	regions, err := client.RegionsList()
	if err != nil {
		return RegionsCache{}, errgo.Notef(err, "fail to list available regions")
	}

	regionsCache.Regions = regions
	regionsCache.ExpireAt = time.Now().Add(10 * time.Minute)

	if opts.SkipAuth {
		return regionsCache, nil
	}

	fd, err = os.OpenFile(c.RegionsCachePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0750)
	if err == nil {
		json.NewEncoder(fd).Encode(regionsCache)
		fd.Close()
	} else {
		debug.Printf("[Regions] Failed to save the regions cache: %v\n", err)
	}

	return regionsCache, nil
}

// GetRegion returns the requested region configuration, use local file system
// cache if any. In case of cache fault, save on disk for 10 minutes the
// available regions
func GetRegion(c Config, name string, opts GetRegionOpts) (scalingo.Region, error) {
	regionsCache, err := EnsureRegionsCache(c, opts)
	if err != nil {
		return scalingo.Region{}, errgo.Notef(err, "fail to get the regions cache")
	}

	for _, region := range regionsCache.Regions {
		if region.Name == name {
			return region, nil
		}
	}
	return scalingo.Region{}, errgo.Notef(err, "invalid region %v", name)
}
