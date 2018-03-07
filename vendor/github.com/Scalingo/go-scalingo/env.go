package scalingo

import "gopkg.in/errgo.v1"

type Variable struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Variables []*Variable

func (vs Variables) Contains(name string) (*Variable, bool) {
	for _, v := range vs {
		if v.Name == name {
			return v, true
		}
	}
	return nil, false
}

type VariablesRes struct {
	Variables Variables `json:"variables"`
}

type VariableSetParams struct {
	Variable *Variable `json:"variable"`
}

func (c *Client) VariablesList(app string) (Variables, error) {
	return c.variableList(app, true)
}

func (c *Client) VariablesListWithoutAlias(app string) (Variables, error) {
	return c.variableList(app, false)
}

func (c *Client) variableList(app string, aliases bool) (Variables, error) {
	var variablesRes VariablesRes
	err := c.subresourceList(app, "variables", map[string]bool{"aliases": aliases}, &variablesRes)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return variablesRes.Variables, nil
}

func (c *Client) VariableSet(app string, name string, value string) (*Variable, int, error) {
	req := &APIRequest{
		Client:   c,
		Method:   "POST",
		Endpoint: "/apps/" + app + "/variables",
		Params: map[string]interface{}{
			"variable": map[string]string{
				"name":  name,
				"value": value,
			},
		},
		Expected: Statuses{200, 201},
	}
	res, err := req.Do()
	if err != nil {
		return nil, 0, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	var params VariableSetParams
	err = ParseJSON(res, &params)
	if err != nil {
		return nil, 0, errgo.Mask(err, errgo.Any)
	}

	return params.Variable, res.StatusCode, nil
}

func (c *Client) VariableMultipleSet(app string, variables Variables) (Variables, int, error) {
	req := &APIRequest{
		Client:   c,
		Method:   "PUT",
		Endpoint: "/apps/" + app + "/variables",
		Params: map[string]Variables{
			"variables": variables,
		},
		Expected: Statuses{200, 201},
	}
	res, err := req.Do()
	if err != nil {
		return nil, 0, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	var params VariablesRes
	err = ParseJSON(res, &params)
	if err != nil {
		return nil, 0, errgo.Mask(err, errgo.Any)
	}

	return params.Variables, res.StatusCode, nil
}

func (c *Client) VariableUnset(app string, id string) error {
	return c.subresourceDelete(app, "variables", id)
}
