package sshkeys

import (
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
	"gopkg.in/errgo.v1"
)

func ReadPrivateKey(path string) (ssh.Signer, error) {
	privateKeyContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return ReadPrivateKeyWithContent(path, privateKeyContent)
}

func ReadPrivateKeyWithContent(path string, privateKeyContent []byte) (ssh.Signer, error) {
	// We parse the private key on our own first so that we can
	// show a nicer error if the private key has a password.
	block, _ := pem.Decode(privateKeyContent)
	if block == nil {
		return nil, fmt.Errorf(
			"Failed to read key '%s': is not in the PEM format", path)
	}

	privateKey := &PrivateKey{Block: block, Path: path}
	if privateKey.IsEncrypted() {
		err := privateKey.Decrypt()
		if err != nil {
			return nil, errgo.Mask(err)
		}
	}

	signer, err := privateKey.Signer()
	if err != nil {
		return nil, errgo.Newf("Invalid SSH key or password: %v", err)
	}

	return signer, nil
}
