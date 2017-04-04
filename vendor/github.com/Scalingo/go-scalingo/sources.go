package scalingo

import "gopkg.in/errgo.v1"

type SourcesCreateResponse struct {
	Source *Source `json:"source"`
}

type Source struct {
	DownloadURL string `json:"download_url"`
	UploadURL   string `json:"upload_url"`
}

func (c *Client) SourcesCreate() (*Source, error) {
	req := &APIRequest{
		Client:   c,
		Method:   "POST",
		Endpoint: "/sources",
		Expected: Statuses{201},
	}
	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer res.Body.Close()

	var sourceRes *SourcesCreateResponse
	err = ParseJSON(res, &sourceRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return sourceRes.Source, nil
}
