package keys

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/crypto/sshkeys"
	"gopkg.in/errgo.v1"
)

func AddAuto() error {
	keyname, err := os.Hostname()
	if err != nil {
		return errgo.Mask(err)
	}

	user, err := user.Current()
	if err != nil {
		return errgo.Mask(err)
	}
	publicKeyPath := filepath.Join(user.HomeDir, ".ssh/id_rsa.pub")

	_, err = os.Stat(publicKeyPath)
	if err != nil {
		privateKeyPath := filepath.Join(user.HomeDir, ".ssh/id_rsa")
		fmt.Println("No key found. Generating one.")

		publicKeyFile, err := os.Create(publicKeyPath)
		if err != nil {
			return errgo.Mask(err)
		}
		defer publicKeyFile.Close()

		privateKeyFile, err := os.Create(privateKeyPath)
		if err != nil {
			return errgo.Mask(err)
		}
		defer privateKeyFile.Close()

		err = privateKeyFile.Chmod(0600)
		if err != nil {
			return errgo.Mask(err)
		}

		publicKey, privateKey, err := sshkeys.GenerateKey()
		if err != nil {
			return errgo.Mask(err)
		}

		_, err = privateKeyFile.Write(privateKey)
		if err != nil {
			return errgo.Mask(err)
		}
		_, err = publicKeyFile.Write(publicKey)
		if err != nil {
			return errgo.Mask(err)
		}
	}

	return Add(keyname, publicKeyPath)
}

func Add(name string, path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return errgo.Mask(err)
	}
	if stat.Mode().IsDir() {
		return errgo.Newf("%s: is a directory", path)
	}
	if stat.Size() > 10*1024*1024 {
		return errgo.Newf("%s: is too large (%v bytes)", stat.Size())
	}

	keyContent, err := ioutil.ReadFile(path)
	if err != nil {
		return errgo.Mask(err)
	}

	c := config.ScalingoClient()
	_, err = c.KeysAdd(name, string(keyContent))
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Printf("Key '%s' has been added.\n", name)
	return nil
}
