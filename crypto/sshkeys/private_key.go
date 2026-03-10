package sshkeys

import (
	"context"
	"encoding/pem"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/Scalingo/cli/term"
	"github.com/Scalingo/go-utils/errors/v3"
)

type PrivateKey struct {
	Path  string
	Block *pem.Block
	PasswordMethod
}

type PasswordMethod func(ctx context.Context, prompt string) (string, error)

func (p *PrivateKey) signer(ctx context.Context) (ssh.Signer, error) {
	if !p.isEncrypted() {
		return ssh.ParsePrivateKey(pem.EncodeToMemory(p.Block))
	}

	if p.PasswordMethod == nil {
		p.PasswordMethod = term.Password
	}

	password, err := p.PasswordMethod(ctx, "Encrypted SSH Key, password: ")
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "fail to get the SSH key password")
	}

	parsedPrivateKey, err := ssh.ParseRawPrivateKeyWithPassphrase(pem.EncodeToMemory(p.Block), []byte(password))
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "parse encrypted private key")
	}

	signer, err := ssh.NewSignerFromKey(parsedPrivateKey)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "not a valid signer")
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
