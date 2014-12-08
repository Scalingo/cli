package keys

import (
	"fmt"

	"github.com/Scalingo/cli/api"
	"gopkg.in/errgo.v1"
)

func Remove(name string) error {
	k, err := keyByName(name)
	if err != nil {
		return errgo.Mask(err)
	}

	err = api.KeysDelete(k.ID)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Printf("Key '%s' has been deleted.\n", name)
	return nil
}

func keyByName(name string) (*api.Key, error) {
	keys, err := api.KeysList()
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
