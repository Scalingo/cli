package api

import "github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"

const (
	EnvNameMaxLength  = 64
	EnvValueMaxLength = 1024
)

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

func VariablesList(app string) (Variables, error) {
	return variableList(app, true)
}

func VariablesListWithoutAlias(app string) (Variables, error) {
	return variableList(app, false)
}

func variableList(app string, aliases bool) (Variables, error) {
	var variablesRes VariablesRes
	err := subresourceList(app, "variables", map[string]bool{"aliases": aliases}, &variablesRes)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return variablesRes.Variables, nil
}

func VariableSet(app string, name string, value string) (*Variable, int, error) {
	req := &APIRequest{
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

func VariableUnset(app string, id string) error {
	return subresourceDelete(app, "variables", id)
}
