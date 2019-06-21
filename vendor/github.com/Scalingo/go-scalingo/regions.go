package scalingo

import (
	"github.com/Scalingo/go-scalingo/http"
	"gopkg.in/errgo.v1"
)

type RegionsService interface {
	RegionsList() ([]Region, error)
}

type Region struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	API         string `json:"api"`
	Dashboard   string `json:"dashboard"`
	DatabaseAPI string `json:"database_api"`
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
