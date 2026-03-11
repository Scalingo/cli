package scalingo

import (
	"context"

	"github.com/Scalingo/go-scalingo/v11/http"
	"github.com/Scalingo/go-utils/errors/v3"
)

type VariablesService interface {
	VariablesList(ctx context.Context, app string) (Variables, error)
	VariablesListWithoutAlias(ctx context.Context, app string) (Variables, error)
	VariableSet(ctx context.Context, app string, name string, value string) (*Variable, error)
	VariableMultipleSet(ctx context.Context, app string, variables Variables) (Variables, error)
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
		return nil, errors.Wrap(ctx, err, "list app variables")
	}
	return variablesRes.Variables, nil
}

func (c *Client) VariableSet(ctx context.Context, app string, name string, value string) (*Variable, error) {
	var variablesSetRes VariableSetParams
	err := c.ScalingoAPI().DoRequest(ctx, &http.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/variables",
		Params: map[string]any{
			"variable": map[string]string{
				"name":  name,
				"value": value,
			},
		},
		Expected: http.Statuses{200, 201},
	}, &variablesSetRes)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "set app variable")
	}

	return variablesSetRes.Variable, nil
}

func (c *Client) VariableMultipleSet(ctx context.Context, app string, variables Variables) (Variables, error) {
	var variabelsRes VariablesRes
	req := &http.APIRequest{
		Method:   "PUT",
		Endpoint: "/apps/" + app + "/variables",
		Params: map[string]Variables{
			"variables": variables,
		},
		Expected: http.Statuses{200, 201},
	}
	err := c.ScalingoAPI().DoRequest(ctx, req, &variabelsRes)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "set multiple app variables")
	}

	return variabelsRes.Variables, nil
}

func (c *Client) VariableUnset(ctx context.Context, app string, id string) error {
	return c.ScalingoAPI().SubresourceDelete(ctx, "apps", app, "variables", id)
}
