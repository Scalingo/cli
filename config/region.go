package config

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
)

type RegionsCache struct {
	ExpireAt time.Time         `json:"expire_at"`
	Regions  []scalingo.Region `json:"regions"`
}

// GetRegion returns the requested region configuration, use local file system cache if any
// In case of cache fault, save on disk for 10 minutes the available regions
func GetRegion(c Config, token, name string) (scalingo.Region, error) {
	var regionsCache RegionsCache
	fd, err := os.Open(c.RegionsCachePath)
	if err == nil {
		json.NewDecoder(fd).Decode(&regionsCache)
		fd.Close()
	}
	if time.Now().Unix() > regionsCache.ExpireAt.Unix() {
		regions, err := ScalingoAuthClientFromToken(token).RegionsList()
		if err != nil {
			return scalingo.Region{}, errgo.Notef(err, "fail to list available regions")
		}
		regionsCache.Regions = regions
		regionsCache.ExpireAt = time.Now().Add(10 * time.Minute)
		fd, err := os.OpenFile(c.RegionsCachePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0750)
		if err == nil {
			json.NewEncoder(fd).Encode(regionsCache)
			fd.Close()
		}
	}
	for _, region := range regionsCache.Regions {
		if region.Name == name {
			return region, nil
		}
	}
	return scalingo.Region{}, errgo.Notef(err, "invalid region %v", name)
}
