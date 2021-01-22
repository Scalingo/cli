package keys

import (
	"fmt"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v4"
	"gopkg.in/errgo.v1"
)

func Remove(name string) error {
	c, err := config.ScalingoAuthClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	k, err := keyByName(c, name)
	if err != nil {
		return errgo.Mask(err)
	}
	err = c.KeysDelete(k.ID)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Printf("Key '%s' has been deleted.\n", name)
	return nil
}

func keyByName(c *scalingo.Client, name string) (*scalingo.Key, error) {
	keys, err := c.KeysList()
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
