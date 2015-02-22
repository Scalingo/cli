package sshkeys

import (
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"strings"

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
		if err, ok := err.(asn1.StructuralError); ok && strings.HasPrefix(err.Msg, "tags don't match") || err.Msg == "length too large" {
			return nil, errgo.Newf("Fail to decrypt SSH key, invalid password.")
		}
		if err, ok := err.(asn1.SyntaxError); ok && err.Msg == "trailing data" {
			return nil, errgo.Newf("The password was OK, but something went wrong.\n" +
				"Please re-run the command with the environment variable DEBUG=1 " +
				"and create an issue with the command output: https://github.com/Scalingo/cli/issues")
		}
		return nil, errgo.Newf("Invalid SSH key or password: %v", err)
	}

	return signer, nil
}
