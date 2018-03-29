package scalingo

import "gopkg.in/errgo.v1"

type KeysService interface {
	KeysList() ([]Key, error)
	KeysAdd(name string, content string) (*Key, error)
	KeysDelete(id string) error
}

var _ KeysService = (*Client)(nil)

type Key struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type KeyIndex struct {
	Keys []Key `json:"keys"`
}

func (c *Client) KeysList() ([]Key, error) {
	req := &APIRequest{
		Client:   c,
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

func (c *Client) KeysAdd(name string, content string) (*Key, error) {
	req := &APIRequest{
		Client:   c,
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

func (c *Client) KeysDelete(id string) error {
	req := &APIRequest{
		Client:   c,
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
