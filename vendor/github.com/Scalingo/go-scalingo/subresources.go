package scalingo

import (
	"gopkg.in/errgo.v1"
)

// subresourceService that wraps the CRUD methods for any subresource of an app on Scalingo.
type subresourceService interface {
	subresourceList(app, subresource string, payload, data interface{}) error
	subresourceAdd(app, subresource string, payload, data interface{}) error
	subresourceGet(app, subresource, id string, payload, data interface{}) error
	subresourceUpdate(app, subresource, id string, payload, data interface{}) error
	subresourceDelete(app, subresource, id string) error
}

type subresourceClient struct {
	*backendConfiguration
}

var _ subresourceService = (*subresourceClient)(nil)

func (c subresourceClient) subresourceGet(app, subresource, id string, payload, data interface{}) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "GET",
		Endpoint: "/apps/" + app + "/" + subresource + "/" + id,
		Params:   payload,
	}, data)
}

func (c subresourceClient) subresourceList(app, subresource string, payload, data interface{}) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "GET",
		Endpoint: "/apps/" + app + "/" + subresource,
		Params:   payload,
	}, data)
}

func (c subresourceClient) subresourceAdd(app, subresource string, payload, data interface{}) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/" + subresource,
		Expected: Statuses{201},
		Params:   payload,
	}, data)
}

func (c subresourceClient) subresourceDelete(app string, subresource string, id string) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "DELETE",
		Endpoint: "/apps/" + app + "/" + subresource + "/" + id,
		Expected: Statuses{204},
	}, nil)
}

func (c subresourceClient) subresourceUpdate(app, subresource, id string, payload, data interface{}) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "PATCH",
		Endpoint: "/apps/" + app + "/" + subresource + "/" + id,
		Params:   payload,
	}, data)
}

func (c subresourceClient) doSubresourceRequest(req *APIRequest, data interface{}) error {
	req.Client = c.backendConfiguration
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
		return errgo.Mask(err, errgo.Any)
	}
	return nil
}
