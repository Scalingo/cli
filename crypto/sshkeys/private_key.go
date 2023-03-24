package sshkeys

import (
	"encoding/pem"
	"strings"

	"github.com/ScaleFT/sshkeys"
	"golang.org/x/crypto/ssh"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/term"
)

type PrivateKey struct {
	Path  string
	Block *pem.Block
	PasswordMethod
}

type PasswordMethod func(prompt string) (string, error)

func (p *PrivateKey) signer() (ssh.Signer, error) {
	if p.IsEncrypted() {
		if p.PasswordMethod == nil {
			p.PasswordMethod = term.Password
		}

		password, err := p.PasswordMethod("Encrypted SSH Key, password: ")
		if err != nil {
			return nil, errgo.Notef(err, "fail to get the SSH key password")
		}

		return sshkeys.ParseEncryptedPrivateKey(pem.EncodeToMemory(p.Block), []byte(password))
	}

	return ssh.ParsePrivateKey(pem.EncodeToMemory(p.Block))
}

func (p *PrivateKey) IsEncrypted() bool {
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
