package keys

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/go-scalingo"
)

func Remove(name string) error {
	k, err := keyByName(name)
	if err != nil {
		return errgo.Mask(err)
	}

	err = scalingo.KeysDelete(k.ID)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Printf("Key '%s' has been deleted.\n", name)
	return nil
}

func keyByName(name string) (*scalingo.Key, error) {
	keys, err := scalingo.KeysList()
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
