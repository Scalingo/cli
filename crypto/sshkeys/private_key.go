package sshkeys

import (
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

var (
	implementedCiphers = []string{"DES-EDE3-CBC", "AES-128-CBC"}
)

type PrivateKey struct {
	Path  string
	Block *pem.Block
	PasswordMethod
}

func (p *PrivateKey) Signer() (ssh.Signer, error) {
	return ssh.ParsePrivateKey(pem.EncodeToMemory(p.Block))
}

func (p *PrivateKey) IsEncrypted() bool {
	return p.Block.Headers["Proc-Type"] == "4,ENCRYPTED"
}

func (p *PrivateKey) IsCipherImplemented(cipher string) bool {
	for _, c := range implementedCiphers {
		if c == cipher {
			return true
		}
	}
	return false
}

func DummyPasswordMethod(password string) PasswordMethod {
	return func(prompt string) (string, error) {
		return prompt, nil
	}
}

type PasswordMethod func(prompt string) (string, error)
