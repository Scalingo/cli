package scalingo

import "gopkg.in/errgo.v1"

type KeysService interface {
	KeysList() ([]Key, error)
	KeysAdd(name string, content string) (*Key, error)
	KeysDelete(id string) error
}

var _ KeysService = (*Client)(nil)

type Key struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type KeyRes struct {
	Key Key `json:"key"`
}

type KeysRes struct {
	Keys []Key `json:"keys"`
}

func (c *Client) KeysList() ([]Key, error) {
	var res KeysRes
	err := c.AuthAPI().ResourceList("keys", nil, &res)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return res.Keys, nil
}

func (c *Client) KeysAdd(name string, content string) (*Key, error) {
	payload := KeyRes{Key{
		Name:    name,
		Content: content,
	}}
	var res KeyRes

	err := c.AuthAPI().ResourceAdd("keys", payload, &res)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return &res.Key, nil
}

func (c *Client) KeysDelete(id string) error {
	err := c.AuthAPI().ResourceDelete("keys", id)
	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}
