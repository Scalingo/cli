package scalingo

import "gopkg.in/errgo.v1"

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

	var source *Source
	err = ParseJSON(res, &source)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return source, nil
}
