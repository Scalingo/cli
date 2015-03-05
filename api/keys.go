package api

import "github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"

type Key struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type KeyIndex struct {
	Keys []Key `json:"keys"`
}

func KeysList() ([]Key, error) {
	req := &APIRequest{
		Endpoint: "/account/keys",
	}
	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer res.Body.Close()

	var ki KeyIndex
	err = ParseJSON(res, &ki)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return ki.Keys, nil
}

func KeysAdd(name string, content string) (*Key, error) {
	req := &APIRequest{
		Method:   "POST",
		Endpoint: "/account/keys",
		Params: map[string]interface{}{
			"key": map[string]interface{}{
				"name":    name,
				"content": content,
			},
		},
		Expected: Statuses{201},
	}
	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer res.Body.Close()

	var key *Key
	err = ParseJSON(res, &key)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return key, nil
}

func KeysDelete(id string) error {
	req := &APIRequest{
		Method:   "DELETE",
		Endpoint: "/account/keys/" + id,
		Expected: Statuses{204},
	}
	res, err := req.Do()
	if err != nil {
		return errgo.Mask(err)
	}
	res.Body.Close()
	return nil
}
