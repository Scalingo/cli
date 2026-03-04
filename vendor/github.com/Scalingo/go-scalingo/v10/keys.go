package scalingo

import (
	"context"

	"github.com/Scalingo/go-utils/errors/v3"
)

type KeysService interface {
	KeysList(context.Context) ([]Key, error)
	KeysAdd(ctx context.Context, name string, content string) (*Key, error)
	KeysDelete(ctx context.Context, id string) error
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

func (c *Client) KeysList(ctx context.Context) ([]Key, error) {
	var res KeysRes
	err := c.AuthAPI().ResourceList(ctx, "keys", nil, &res)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list keys")
	}
	return res.Keys, nil
}

func (c *Client) KeysAdd(ctx context.Context, name string, content string) (*Key, error) {
	payload := KeyRes{Key{
		Name:    name,
		Content: content,
	}}
	var res KeyRes

	err := c.AuthAPI().ResourceAdd(ctx, "keys", payload, &res)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "add key")
	}

	return &res.Key, nil
}

func (c *Client) KeysDelete(ctx context.Context, id string) error {
	err := c.AuthAPI().ResourceDelete(ctx, "keys", id)
	if err != nil {
		return errors.Wrap(ctx, err, "delete key")
	}
	return nil
}
