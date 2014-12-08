package api

import "gopkg.in/errgo.v1"

type Key struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type KeyIndex struct {
	Keys []Key `json:"keys"`
}

func KeysList() ([]Key, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"endpoint": "/keys",
		"expected": Statuses{200},
	}
	res, err := Do(req)
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
	req := map[string]interface{}{
		"method":   "POST",
		"endpoint": "/keys",
		"params": map[string]interface{}{
			"key": map[string]interface{}{
				"name":    name,
				"content": content,
			},
		},
		"expected": Statuses{201},
	}
	res, err := Do(req)
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
	req := map[string]interface{}{
		"method":   "DELETE",
		"endpoint": "/keys/" + id,
		"expected": Statuses{204},
	}
	res, err := Do(req)
	if err != nil {
		return errgo.Mask(err)
	}
	res.Body.Close()
	return nil
}
