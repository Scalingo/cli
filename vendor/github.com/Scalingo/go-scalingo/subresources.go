package scalingo

import (
	"gopkg.in/errgo.v1"
)

// SubresourceService that wraps the CRUD methods for any subresource of an app on Scalingo.
type SubresourceService interface {
	subresourceList(app, subresource string, payload, data interface{}) error
	subresourceAdd(app, subresource string, payload, data interface{}) error
	subresourceGet(app, subresource, id string, payload, data interface{}) error
	subresourceUpdate(app, subresource, id string, payload, data interface{}) error
	subresourceDelete(app, subresource, id string) error
}

var _ SubresourceService = (*Client)(nil)

func (c *Client) subresourceGet(app, subresource, id string, payload, data interface{}) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "GET",
		Endpoint: "/apps/" + app + "/" + subresource + "/" + id,
		Params:   payload,
	}, data)
}

func (c *Client) subresourceList(app, subresource string, payload, data interface{}) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "GET",
		Endpoint: "/apps/" + app + "/" + subresource,
		Params:   payload,
	}, data)
}

func (c *Client) subresourceAdd(app, subresource string, payload, data interface{}) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/" + subresource,
		Expected: Statuses{201},
		Params:   payload,
	}, data)
}

func (c *Client) subresourceDelete(app string, subresource string, id string) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "DELETE",
		Endpoint: "/apps/" + app + "/" + subresource + "/" + id,
		Expected: Statuses{204},
	}, nil)
}

func (c *Client) subresourceUpdate(app, subresource, id string, payload, data interface{}) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "PATCH",
		Endpoint: "/apps/" + app + "/" + subresource + "/" + id,
		Params:   payload,
	}, data)
}

func (c *Client) doSubresourceRequest(req *APIRequest, data interface{}) error {
	req.Client = c
	res, err := req.Do()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	if data == nil {
		return nil
	}

	err = ParseJSON(res, data)
	if err != nil {
		return errgo.NoteMask(err, "fail to parse json of subresource response")
	}
	return nil
}
