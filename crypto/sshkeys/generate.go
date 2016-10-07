package sshkeys

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
	"gopkg.in/errgo.v1"
)

func GenerateKey() ([]byte, []byte, error) {
	var publicKeyBuffer bytes.Buffer
	var privateKeyBuffer bytes.Buffer

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)

	if err != nil {
		return nil, nil, errgo.Mask(err)
	}

	err = privateKey.Validate()
	if err != nil {
		return nil, nil, errgo.Mask(err)
	}

	privateKeyDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privatePEM := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDER,
	}

	err = pem.Encode(&privateKeyBuffer, &privatePEM)
	if err != nil {
		return nil, nil, errgo.Mask(err)
	}

	publicKey := privateKey.PublicKey
	publicKeySSH, err := ssh.NewPublicKey(&publicKey)

	if err != nil {
		return nil, nil, errgo.Mask(err)
	}

	publicKeyBuffer.Write(ssh.MarshalAuthorizedKey(publicKeySSH))
	err = AddKeyToAgent(privateKey)
	if err != nil {
		return nil, nil, errgo.Mask(err)
	}
	return publicKeyBuffer.Bytes(), privateKeyBuffer.Bytes(), nil
}
