package scalingo

import "gopkg.in/errgo.v1"

type SourcesService interface {
	SourcesCreate() (*Source, error)
}

var _ SourcesService = (*Client)(nil)

type SourcesCreateResponse struct {
	Source *Source `json:"source"`
}

type Source struct {
	DownloadURL string `json:"download_url"`
	UploadURL   string `json:"upload_url"`
}

func (c *Client) SourcesCreate() (*Source, error) {
	var sourceRes SourcesCreateResponse
	err := c.ScalingoAPI().ResourceAdd("sources", nil, &sourceRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return sourceRes.Source, nil
}
