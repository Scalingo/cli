package keys

import (
	"context"
	"fmt"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v10"
	"github.com/Scalingo/go-utils/errors/v3"
)

func Remove(ctx context.Context, name string) error {
	c, err := config.ScalingoAuthClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	k, err := keyByName(ctx, c, name)
	if err != nil {
		return errors.Wrapf(ctx, err, "find SSH key %s", name)
	}
	err = c.KeysDelete(ctx, k.ID)
	if err != nil {
		return errors.Wrapf(ctx, err, "delete SSH key %s", name)
	}

	fmt.Printf("Key '%s' has been deleted.\n", name)
	return nil
}

func keyByName(ctx context.Context, c *scalingo.Client, name string) (*scalingo.Key, error) {
	keys, err := c.KeysList(ctx)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list SSH keys")
	}
	for _, k := range keys {
		if k.Name == name {
			return &k, nil
		}
	}
	return nil, errors.New(ctx, "no such key")
}
