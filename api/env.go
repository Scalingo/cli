package api

import "gopkg.in/errgo.v1"

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

type VariablesListParams struct {
	Variables Variables `json:"variables"`
}

type VariableSetParams struct {
	Variable *Variable `json:"variable"`
}

func VariablesList(app string) (Variables, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"endpoint": "/apps/" + app + "/variables",
		"expected": Statuses{200},
	}
	res, err := Do(req)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	var params VariablesListParams
	err = ParseJSON(res, &params)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return params.Variables, nil
}

func VariableSet(app string, name string, value string) (*Variable, int, error) {
	req := map[string]interface{}{
		"method":   "POST",
		"endpoint": "/apps/" + app + "/variables",
		"params": map[string]interface{}{
			"variable": map[string]string{
				"name":  name,
				"value": value,
			},
		},
		"expected": Statuses{200, 201},
	}
	res, err := Do(req)
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
	req := map[string]interface{}{
		"method":    "DELETE",
		"endpoint":  "/apps/" + app + "/variables/" + id,
		"exepected": Statuses{204},
	}
	_, err := Do(req)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	return nil
}
