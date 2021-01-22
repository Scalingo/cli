package scalingo

import (
	"sync"

	"github.com/Scalingo/go-scalingo/v4/http"
	"gopkg.in/errgo.v1"
)

var (
	ErrRegionNotFound = errgo.New("Region not found")

	regionCache      map[string]Region = map[string]Region{}
	regionCacheMutex *sync.Mutex       = &sync.Mutex{}
)

type RegionsService interface {
	RegionsList() ([]Region, error)
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

func (c *Client) RegionsList() ([]Region, error) {
	var res regionsRes
	err := c.AuthAPI().DoRequest(&http.APIRequest{
		Method:   "GET",
		Endpoint: "/regions",
	}, &res)
	if err != nil {
		return nil, errgo.Notef(err, "fail to call GET /regions")
	}
	return res.Regions, nil
}

func (c *Client) getRegion(regionName string) (Region, error) {
	regionCacheMutex.Lock()
	defer regionCacheMutex.Unlock()

	if _, ok := regionCache[regionName]; !ok {
		regions, err := c.RegionsList()
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
