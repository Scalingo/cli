package keys

import (
	"context"
	"fmt"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v7"
)

func Remove(ctx context.Context, name string) error {
	c, err := config.ScalingoAuthClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	k, err := keyByName(ctx, c, name)
	if err != nil {
		return errgo.Mask(err)
	}
	err = c.KeysDelete(ctx, k.ID)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Printf("Key '%s' has been deleted.\n", name)
	return nil
}

func keyByName(ctx context.Context, c *scalingo.Client, name string) (*scalingo.Key, error) {
	keys, err := c.KeysList(ctx)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	for _, k := range keys {
		if k.Name == name {
			return &k, nil
		}
	}
	return nil, errgo.New("no such key")
}
