package keys

import (
	"fmt"
	"os"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func Add(name string, path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return errgo.Notef(err, "fail to stat path to the key")
	}
	if stat.Mode().IsDir() {
		return fmt.Errorf("%s: is a directory", path)
	}
	if stat.Size() > 10*1024*1024 {
		return fmt.Errorf("%s: is too large (%v bytes)", path, stat.Size())
	}

	keyContent, err := os.ReadFile(path)
	if err != nil {
		return errgo.Notef(err, "fail to read the key file")
	}

	c, err := config.ScalingoAuthClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	_, err = c.KeysAdd(name, string(keyContent))
	if err != nil {
		return errgo.Notef(err, "fail to add the key to Scalingo account")
	}

	fmt.Printf("Key '%s' has been added.\n", name)
	return nil
}
