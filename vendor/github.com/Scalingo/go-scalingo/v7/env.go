package scalingo

import (
	"context"
	"encoding/json"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v7/http"
)

type VariablesService interface {
	VariablesList(ctx context.Context, app string) (Variables, error)
	VariablesListWithoutAlias(ctx context.Context, app string) (Variables, error)
	VariableSet(ctx context.Context, app string, name string, value string) (*Variable, int, error)
	VariableMultipleSet(ctx context.Context, app string, variables Variables) (Variables, int, error)
	VariableUnset(ctx context.Context, app string, id string) error
}

var _ VariablesService = (*Client)(nil)

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

func (c *Client) VariablesList(ctx context.Context, app string) (Variables, error) {
	return c.variableList(ctx, app, true)
}

func (c *Client) VariablesListWithoutAlias(ctx context.Context, app string) (Variables, error) {
	return c.variableList(ctx, app, false)
}

func (c *Client) variableList(ctx context.Context, app string, aliases bool) (Variables, error) {
	var variablesRes VariablesRes
	err := c.ScalingoAPI().SubresourceList(ctx, "apps", app, "variables", map[string]bool{"aliases": aliases}, &variablesRes)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return variablesRes.Variables, nil
}

func (c *Client) VariableSet(ctx context.Context, app string, name string, value string) (*Variable, int, error) {
	req := &http.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/variables",
		Params: map[string]interface{}{
			"variable": map[string]string{
				"name":  name,
				"value": value,
			},
		},
		Expected: http.Statuses{200, 201},
	}
	res, err := c.ScalingoAPI().Do(ctx, req)
	if err != nil {
		return nil, 0, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	var params VariableSetParams
	err = json.NewDecoder(res.Body).Decode(&params)
	if err != nil {
		return nil, 0, errgo.Mask(err, errgo.Any)
	}

	return params.Variable, res.StatusCode, nil
}

func (c *Client) VariableMultipleSet(ctx context.Context, app string, variables Variables) (Variables, int, error) {
	req := &http.APIRequest{
		Method:   "PUT",
		Endpoint: "/apps/" + app + "/variables",
		Params: map[string]Variables{
			"variables": variables,
		},
		Expected: http.Statuses{200, 201},
	}
	res, err := c.ScalingoAPI().Do(ctx, req)
	if err != nil {
		return nil, 0, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	var params VariablesRes
	err = json.NewDecoder(res.Body).Decode(&params)
	if err != nil {
		return nil, 0, errgo.Mask(err, errgo.Any)
	}

	return params.Variables, res.StatusCode, nil
}

func (c *Client) VariableUnset(ctx context.Context, app string, id string) error {
	return c.ScalingoAPI().SubresourceDelete(ctx, "apps", app, "variables", id)
}
