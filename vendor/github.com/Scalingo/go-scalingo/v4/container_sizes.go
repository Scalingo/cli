package scalingo

import (
	"gopkg.in/errgo.v1"

	httpclient "github.com/Scalingo/go-scalingo/v4/http"
)

type ContainerSize struct {
	ID        string `json:"id"`
	SKU       string `json:"sku,omitempty"`
	Name      string `json:"name"`
	HumanName string `json:"human_name"`
	HumanCPU  string `json:"human_cpu"`
	Memory    int    `json:"memory"`
	PidsLimit int    `json:"pids_limit,omitempty"`

	HourlyPrice     int                          `json:"hourly_price"`
	ThirtydaysPrice int                          `json:"thirtydays_price"`
	Pricings        map[string]map[string]string `json:"pricings"`

	Swap    int `json:"swap"`
	Ordinal int `json:"ordinal"`
}

type ContainerSizesService interface {
	ContainerSizesList() ([]ContainerSize, error)
}

var _ ContainerSizesService = (*Client)(nil)

func (c *Client) ContainerSizesList() ([]ContainerSize, error) {
	req := &httpclient.APIRequest{
		Endpoint: "/features/container_sizes",
	}

	resmap := map[string][]ContainerSize{}
	err := c.ScalingoAPI().DoRequest(req, &resmap)
	if err != nil {
		return nil, errgo.Notef(err, "fail to request Scalingo API to list the container sizes")
	}
	return resmap["container_sizes"], nil
}
