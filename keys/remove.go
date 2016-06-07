package keys

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/config"
)

func Remove(name string) error {
	k, err := keyByName(name)
	if err != nil {
		return errgo.Mask(err)
	}

	c := config.ScalingoClient()
	err = c.KeysDelete(k.ID)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Printf("Key '%s' has been deleted.\n", name)
	return nil
}

func keyByName(name string) (*scalingo.Key, error) {
	c := config.ScalingoClient()
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
