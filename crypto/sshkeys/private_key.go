package sshkeys

import (
	"context"
	"encoding/pem"
	"strings"

	"golang.org/x/crypto/ssh"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/term"
	errors "github.com/Scalingo/go-utils/errors/v2"
)

type PrivateKey struct {
	Path  string
	Block *pem.Block
	PasswordMethod
}

type PasswordMethod func(prompt string) (string, error)

func (p *PrivateKey) signer() (ssh.Signer, error) {
	ctx := context.TODO()

	if !p.isEncrypted() {
		return ssh.ParsePrivateKey(pem.EncodeToMemory(p.Block))
	}

	if p.PasswordMethod == nil {
		p.PasswordMethod = term.Password
	}

	password, err := p.PasswordMethod("Encrypted SSH Key, password: ")
	if err != nil {
		return nil, errgo.Notef(err, "fail to get the SSH key password")
	}

	parsedPrivateKey, err := ssh.ParseRawPrivateKeyWithPassphrase(pem.EncodeToMemory(p.Block), []byte(password))
	if err != nil {
		return nil, errors.Notef(ctx, err, "parse encrypted private key")
	}

	signer, ok := parsedPrivateKey.(ssh.Signer)
	if !ok {
		// ssh.ParseRawPrivateKeyWithPassphrase returns an empty interface for
		// retro-compatibility reasons, all private key types in the standard library implement
		// [...] interfaces such as Signer.
		// Hence this error should never happen.
		// https://pkg.go.dev/crypto@go1.20.2#PrivateKey
		return nil, errors.New(ctx, "not a valid signer")
	}
	return signer, nil
}

func (p *PrivateKey) isEncrypted() bool {
	return p.Block.Headers["Proc-Type"] == "4,ENCRYPTED" || p.isOpenSSHFormatEncrypted()
}

func (p *PrivateKey) isOpenSSHFormatEncrypted() bool {
	if p.Block.Type != "OPENSSH PRIVATE KEY" {
		return false
	}

	_, err := ssh.ParseRawPrivateKey(pem.EncodeToMemory(p.Block))
	if err != nil {
		return strings.Contains(err.Error(), "this private key is passphrase protected")
	}
	return false
}
