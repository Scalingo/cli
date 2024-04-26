package scalingo

import (
	"context"
	"sync"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v7/http"
)

var (
	ErrRegionNotFound = errgo.New("Region not found")

	regionCache      = map[string]Region{}
	regionCacheMutex = &sync.Mutex{}
)

type RegionsService interface {
	RegionsList(context.Context) ([]Region, error)
}

type Region struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	SSH         string `json:"ssh"`
	API         string `json:"api"`
	Dashboard   string `json:"dashboard"`
	DatabaseAPI string `json:"database_api"`
	Default     bool   `json:"default"`
}

type regionsRes struct {
	Regions []Region `json:"regions"`
}

func (c *Client) RegionsList(ctx context.Context) ([]Region, error) {
	var res regionsRes
	err := c.AuthAPI().DoRequest(ctx, &http.APIRequest{
		Method:   "GET",
		Endpoint: "/regions",
	}, &res)
	if err != nil {
		return nil, errgo.Notef(err, "fail to call GET /regions")
	}
	return res.Regions, nil
}

func (c *Client) getRegion(ctx context.Context, regionName string) (Region, error) {
	regionCacheMutex.Lock()
	defer regionCacheMutex.Unlock()

	if _, ok := regionCache[regionName]; !ok {
		regions, err := c.RegionsList(ctx)
		if err != nil {
			return Region{}, errgo.Notef(err, "fail to list regions")
		}

		for _, region := range regions {
			regionCache[region.Name] = region
		}
	}

	region, ok := regionCache[regionName]
	if !ok {
		return Region{}, ErrRegionNotFound
	}
	return region, nil
}
