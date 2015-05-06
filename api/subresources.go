package api

import "github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"

func subresourceList(app, subresource string, payload, data interface{}) error {
	req := &APIRequest{
		Endpoint: "/apps/" + app + "/" + subresource,
		Params:   payload,
	}
	res, err := req.Do()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	err = ParseJSON(res, data)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	return nil
}

func subresourceAdd(app, subresource string, payload, data interface{}) error {
	req := &APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/" + subresource,
		Expected: Statuses{201},
		Params:   payload,
	}

	res, err := req.Do()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	err = ParseJSON(res, data)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	return nil
}

func subresourceDelete(app string, subresource string, id string) error {
	req := &APIRequest{
		Method:   "DELETE",
		Endpoint: "/apps/" + app + "/" + subresource + "/" + id,
		Expected: Statuses{204},
	}

	res, err := req.Do()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	res.Body.Close()
	return nil
}

func subresourceUpdate(app, subresource, id string, payload, data interface{}) error {
	req := &APIRequest{
		Method:   "PATCH",
		Endpoint: "/apps/" + app + "/" + subresource + "/" + id,
		Params:   payload,
	}

	res, err := req.Do()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	err = ParseJSON(res, data)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	return nil
}
